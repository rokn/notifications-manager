package main

import (
	"github.com/rokn/notifications-manager/pkg/channels"
	"github.com/rokn/notifications-manager/pkg/config"
	"github.com/rokn/notifications-manager/pkg/notifiers"
	"github.com/rokn/notifications-manager/pkg/queue"
	"go.uber.org/zap"
)

func initNotifier(notifierType channels.ChannelType, logger *zap.Logger) notifiers.Notifier {
	switch notifierType {
	case channels.Email:
		return notifiers.NewEmailNotifier(logger)
	case channels.Slack:
		return notifiers.NewSlackNotifier(logger)
	default:
		logger.Fatal("unsupported notifier type", zap.String("type", string(notifierType)))
	}
	return nil
}

func main() {
	cfg := notifierConfig{}
	logger := config.InitConfigWithLogger(&cfg)

	logger.Info("connecting to channels server", zap.String("url", cfg.ChannelsServerURL))
	channelsClient := channels.NewClient(cfg.ChannelsServerURL, logger)
	defer func() {
		err := channelsClient.Close()
		if err != nil {
			logger.Fatal("failed to close channels client", zap.Error(err))
		}
	}()

	logger.Info("initializing notifier", zap.String("type", string(cfg.ChannelType)))
	notifier := initNotifier(cfg.ChannelType, logger)
	handler := notifiers.NewNotificationHandler(channelsClient, notifier, logger)

	logger.Info("starting consumer", zap.String("uri", cfg.RabbitURI))
	consumer := queue.NewRabbitConsumer(cfg.RabbitURI, cfg.ChannelType, queue.DefaultNotificationRouter(), handler, logger)
	if err := consumer.Start(); err != nil {
		logger.Fatal("failed to start consumer", zap.Error(err))
	}
}
