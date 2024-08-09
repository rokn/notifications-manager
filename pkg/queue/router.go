package queue

import (
	"fmt"
	"github.com/rokn/notifications-manager/pkg/channels"
	"strings"
)

type NotificationRouter interface {
	GetRoutingKey(channelType channels.ChannelType, channel string) string
	GetRoutingKeyForType(channelType channels.ChannelType) string
	GetQueueName(channelType channels.ChannelType) string
	GetChannelName(routingKey string) (string, error)
	GetExchange() string
}

type defaultNotificationRouter struct{}

func DefaultNotificationRouter() NotificationRouter {
	return &defaultNotificationRouter{}
}

func (d *defaultNotificationRouter) GetQueueName(channelType channels.ChannelType) string {
	return fmt.Sprintf("%s.%s", d.GetExchange(), channelType)
}

func (d *defaultNotificationRouter) GetRoutingKeyForType(channelType channels.ChannelType) string {
	return fmt.Sprintf("%s.%s.*", d.GetExchange(), channelType)
}

func (d *defaultNotificationRouter) GetRoutingKey(channelType channels.ChannelType, channel string) string {
	return fmt.Sprintf("%s.%s.%s", d.GetExchange(), channelType, channel)
}

func (d *defaultNotificationRouter) GetChannelName(routingKey string) (string, error) {
	parts := strings.Split(routingKey, ".")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid routing key: %s", routingKey)
	}
	return parts[2], nil
}

func (d *defaultNotificationRouter) GetExchange() string {
	return "notifications"
}
