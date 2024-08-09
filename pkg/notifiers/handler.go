package notifiers

import (
	"context"
	"github.com/rokn/notifications-manager/pkg/channels"
	"github.com/rokn/notifications-manager/pkg/queue"
	"go.uber.org/zap"
)

type Notifier interface {
	SendNotification(channel *channels.ChannelModel, notification queue.NotificationMessage) error
	SupportedChannelTypes() []channels.ChannelType
}

type defaultHandler struct {
	channelsClient channels.Client
	notifier       Notifier
	log            *zap.Logger
}

func NewNotificationHandler(channelsClient channels.Client, notifier Notifier, logger *zap.Logger) queue.NotificationHandler {
	return &defaultHandler{
		channelsClient: channelsClient,
		notifier:       notifier,
		log:            logger.With(zap.String("service", "handler")),
	}
}

func (s *defaultHandler) HandleNotification(channel string, notification queue.NotificationMessage) error {
	// Get the channel configuration
	channelInfo, err := s.channelsClient.GetChannel(context.Background(), channel)
	if err != nil {
		s.log.Error("failed to get channel configuration", zap.Error(err))
		return err
	}
	supportedChannels := s.notifier.SupportedChannelTypes()
	if !channelInfo.Type.In(supportedChannels) {
		s.log.Error("channel type not supported", zap.String("channel", channel), zap.Any("supported", supportedChannels))
		return nil
	}

	return s.notifier.SendNotification(channelInfo, notification)
}
