package main

import (
	"github.com/rokn/notifications-manager/pkg/channels"
	"github.com/rokn/notifications-manager/pkg/config"
	"github.com/rokn/notifications-manager/pkg/ingress"
	"github.com/rokn/notifications-manager/pkg/queue"
	"go.uber.org/zap"
)

func main() {
	cfg := ingressConfig{}
	logger := config.InitConfigWithLogger(&cfg)

	channelsClient := channels.NewClient(cfg.ChannelsServerURL, logger)
	defer func() {
		err := channelsClient.Close()
		if err != nil {
			logger.Fatal("failed to close channels client", zap.Error(err))
		}
	}()

	queuePublisher := queue.NewRabbitPublisher(cfg.RabbitURI, queue.DefaultNotificationRouter(), logger)
	defer func() {
		err := queuePublisher.Close()
		if err != nil {
			logger.Fatal("failed to close queue publisher", zap.Error(err))
		}
	}()

	service := ingress.NewService(channelsClient, queuePublisher, logger)
	server := ingress.NewServer(cfg.Port, cfg.Profile, service, logger)
	if err := server.Start(); err != nil {
		logger.Fatal("server exited with error", zap.Error(err))
	}
}
