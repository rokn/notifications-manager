package channels

import (
	"context"
	"encoding/json"
	"github.com/rokn/notifications-manager/pkg/channels/api"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	api.UnimplementedChannelsServiceServer
	svc Service
	log *zap.Logger
}

// NewServer creates a new channels server.
func NewServer(svc Service, logger *zap.Logger) api.ChannelsServiceServer {
	return &server{
		svc: svc,
		log: logger.With(zap.String("server", "channels")),
	}
}

// GetChannel returns the configuration for a channel.
func (s *server) GetChannel(_ context.Context, req *api.GetChannelRequest) (*api.GetChannelResponse, error) {
	channel, err := s.svc.GetChannel(req.Name)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	jsonConfig, err := json.Marshal(channel.Configuration)
	if err != nil {
		s.log.Warn("failed to encode configuration", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to encode configuration")
	}

	return &api.GetChannelResponse{
		Name:          *channel.Name,
		Type:          string(*channel.Type),
		Configuration: jsonConfig,
	}, nil
}

// GetChannelNames returns the names of all channels.
func (s *server) GetChannelNames(_ context.Context, _ *api.GetChannelNamesRequest) (*api.GetChannelNamesResponse, error) {
	names := s.svc.GetChannelNames()
	return &api.GetChannelNamesResponse{Names: names}, nil
}
