package queue

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rokn/notifications-manager/pkg/channels"
	"go.uber.org/zap"
)

type NotificationHandler interface {
	HandleNotification(channel string, n NotificationMessage) error
}

type Consumer interface {
	Start() error
}

type rabbitConsumer struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	channelType channels.ChannelType
	router      NotificationRouter
	handler     NotificationHandler
	log         *zap.Logger
}

func NewRabbitConsumer(
	rabbitUri string,
	channelType channels.ChannelType,
	notificationRouter NotificationRouter,
	handler NotificationHandler,
	logger *zap.Logger,
) Consumer {
	log := logger.With(zap.String("service", "rabbit_consumer"))
	conn, err := amqp.Dial(rabbitUri)
	if err != nil {
		log.Fatal("failed to connect to RabbitMQ", zap.Error(err))
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("failed to open a channel", zap.Error(err))
	}

	err = ch.ExchangeDeclare(
		notificationRouter.GetExchange(),
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)

	return &rabbitConsumer{
		conn:        conn,
		ch:          ch,
		channelType: channelType,
		log:         log,
		handler:     handler,
		router:      notificationRouter,
	}
}

func (r *rabbitConsumer) Start() error {
	defer func() {
		err := r.conn.Close()
		if err != nil {
			r.log.Error("failed to close connection", zap.Error(err))
		}
	}()

	q, err := r.ch.QueueDeclare(
		r.router.GetQueueName(r.channelType),
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		r.log.Error("failed to declare a queue", zap.Error(err))
		return err
	}

	err = r.ch.QueueBind(
		q.Name, // queue name
		r.router.GetRoutingKeyForType(r.channelType),
		r.router.GetExchange(),
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		r.log.Error("failed to bind queue", zap.Error(err))
		return err
	}

	msgs, err := r.ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // arguments
	)

	for msg := range msgs {
		r.log.Info("received message", zap.String("routing_key", msg.RoutingKey))
		channel, parsedBody, err := r.parseMessage(&msg)
		if err != nil {
			// Nack the message if we can't parse it
			r.nackMessage(&msg, false)
			continue
		}

		err = r.handler.HandleNotification(*channel, *parsedBody)
		if err != nil {
			r.log.Error("failed to handle notification", zap.Error(err))
			r.nackMessage(&msg, true)
			continue
		}

		err = msg.Ack(false)
		if err != nil {
			r.log.Error("failed to ack message", zap.Error(err))
		}
	}

	return nil
}

func (r *rabbitConsumer) nackMessage(msg *amqp.Delivery, requeue bool) {
	err := msg.Nack(false, false)
	if err != nil {
		r.log.Error("failed to nack message", zap.Error(err))
	}
}

func (r *rabbitConsumer) parseMessage(msg *amqp.Delivery) (*string, *NotificationMessage, error) {
	channel, err := r.router.GetChannelName(msg.RoutingKey)
	if err != nil {
		r.log.Error("failed to get channel from routing key", zap.Error(err))
		return nil, nil, err
	}

	parsedBody := NotificationMessage{}
	err = json.Unmarshal(msg.Body, &parsedBody)
	if err != nil {
		r.log.Error("failed to unmarshal message body", zap.Error(err))
		return nil, nil, err
	}

	return &channel, &parsedBody, nil
}
