package main

import (
	"fmt"
	"github.com/rokn/notifications-manager/pkg/channels"
	"github.com/rokn/notifications-manager/pkg/channels/api"
	"github.com/rokn/notifications-manager/pkg/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := channelsConfig{}
	logger := config.InitConfigWithLogger(&cfg)
	logger.Info("starting channels service", zap.Int("port", cfg.Port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	logger.Info("initializing grpc server")
	s := grpc.NewServer()
	svc := channels.NewService(cfg.ChannelsConfig, logger)
	api.RegisterChannelsServiceServer(s, channels.NewServer(svc, logger))
	logger.Info("starting grpc server")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
