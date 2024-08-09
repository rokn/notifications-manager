package notifiers

import (
	"github.com/rokn/notifications-manager/pkg/channels"
	"github.com/rokn/notifications-manager/pkg/queue"
	"go.uber.org/zap"
)

type emailNotifier struct {
	log *zap.Logger
}

func NewEmailNotifier(logger *zap.Logger) Notifier {
	return &emailNotifier{
		log: logger.With(zap.String("notifier", "email")),
	}
}

func (s *emailNotifier) SendNotification(channel *channels.ChannelModel, notification queue.NotificationMessage) error {
	s.log.Info("sending email notification", zap.String("channel", channel.Name), zap.Any("notification", notification))
	return nil
}

func (s *emailNotifier) SupportedChannelTypes() []channels.ChannelType {
	return []channels.ChannelType{channels.Email}
}
