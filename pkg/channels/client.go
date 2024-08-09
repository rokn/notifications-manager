package channels

import (
	"context"
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"github.com/rokn/notifications-manager/pkg/channels/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const ChannelNamesCacheKey = "channel_names"

type Client interface {
	GetChannel(ctx context.Context, name string) (*ChannelModel, error)
	GetChannelNames(ctx context.Context) ([]string, error)
	Close() error
}

type client struct {
	c     api.ChannelsServiceClient
	cache *cache.Cache
	conn  *grpc.ClientConn
	log   *zap.Logger
}

func getChannelCacheKey(name string) string {
	return "channel:" + name
}

func NewClient(serverUrl string, logger *zap.Logger) Client {
	log := logger.With(zap.String("client", "channels"))
	conn, err := grpc.NewClient(serverUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to channels server", zap.Error(err))
	}

	return &client{
		conn:  conn,
		c:     api.NewChannelsServiceClient(conn),
		cache: cache.New(5*time.Minute, 10*time.Minute),
		log:   log,
	}
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) GetChannel(ctx context.Context, name string) (*ChannelModel, error) {
	if channel, found := c.cache.Get(getChannelCacheKey(name)); found {
		return channel.(*ChannelModel), nil
	}

	resp, err := c.c.GetChannel(ctx, &api.GetChannelRequest{Name: name})
	if err != nil {
		return nil, err
	}

	config := make(map[string]string)
	err = json.Unmarshal(resp.Configuration, &config)
	if err != nil {
		return nil, err
	}

	result := &ChannelModel{
		Name:          resp.Name,
		Type:          ChannelType(resp.Type),
		Configuration: config,
	}

	c.cache.Set(name, result, cache.DefaultExpiration)
	return result, nil
}

func (c *client) GetChannelNames(ctx context.Context) ([]string, error) {
	if names, found := c.cache.Get(ChannelNamesCacheKey); found {
		return names.([]string), nil
	}

	resp, err := c.c.GetChannelNames(ctx, &api.GetChannelNamesRequest{})
	if err != nil {
		return nil, err
	}

	c.cache.Set(ChannelNamesCacheKey, resp.Names, cache.DefaultExpiration)
	return resp.Names, nil
}
