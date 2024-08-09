package queue

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Publisher interface {
	PublishNotification(request PublishRequest) error
	Close() error
}

type rabbitPublisher struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	log    *zap.Logger
	router NotificationRouter
}

func NewRabbitPublisher(rabbitUri string, notificationRouter NotificationRouter, logger *zap.Logger) Publisher {
	log := logger.With(zap.String("service", "rabbit_publisher"))
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
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	return &rabbitPublisher{
		conn:   conn,
		ch:     ch,
		log:    log,
		router: notificationRouter,
	}

}

func (p *rabbitPublisher) PublishNotification(request PublishRequest) error {
	message, err := json.Marshal(request.Message)
	if err != nil {
		p.log.Error("failed to marshal message", zap.Error(err))
		return err
	}

	err = p.ch.Tx() // Start transaction
	if err != nil {
		p.log.Error("failed to start transaction", zap.Error(err))
		return err
	}

	for chType, channels := range request.Channels {
		for _, channel := range channels {
			err = p.ch.Publish(
				p.router.GetExchange(),
				p.router.GetRoutingKey(chType, channel),
				true,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        message,
				})
			if err != nil {
				if txErr := p.ch.TxRollback(); txErr != nil {
					return txErr
				}
				p.log.Error("failed to publish message", zap.Error(err))
				return err
			}
		}
	}

	err = p.ch.TxCommit() // Commit the transaction
	if err != nil {
		p.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	p.log.Debug("message published successfully")
	return nil
}

func (p *rabbitPublisher) Close() error {
	return p.conn.Close()
}
