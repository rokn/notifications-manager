package ingress

import (
	"context"
	"errors"
	"github.com/rokn/notifications-manager/pkg/channels"
	"github.com/rokn/notifications-manager/pkg/queue"
	"go.uber.org/zap"
)

type Service interface {
	TransmitNotification(ctx context.Context, notification NotificationDTO) error
	GetChannels(ctx context.Context) ([]string, error)
}

type service struct {
	channelsClient channels.Client
	queuePublisher queue.Publisher
	log            *zap.Logger
}

func NewService(channelsClient channels.Client, publisher queue.Publisher, logger *zap.Logger) Service {
	return &service{
		channelsClient: channelsClient,
		queuePublisher: publisher,
		log:            logger.With(zap.String("service", "ingress")),
	}
}

func (s *service) TransmitNotification(ctx context.Context, notification NotificationDTO) error {
	s.log.Debug("received notification", zap.Strings("channels", notification.Channels))
	channelsMap := make(map[channels.ChannelType][]string)
	for _, channelName := range notification.Channels {
		channel, err := s.channelsClient.GetChannel(ctx, channelName)
		if err != nil {
			s.log.Info("channel not found", zap.String("channel", channelName))
			return errors.New("non-existent channel")
		}
		channelsMap[channel.Type] = append(channelsMap[channel.Type], channel.Name)
	}

	return s.queuePublisher.PublishNotification(queue.PublishRequest{
		Channels: channelsMap,
		Message: queue.NotificationMessage{
			Title:    notification.Title,
			Message:  notification.Message,
			Metadata: notification.Metadata,
		},
	})
}

func (s *service) GetChannels(ctx context.Context) ([]string, error) {
	return s.channelsClient.GetChannelNames(ctx)
}
