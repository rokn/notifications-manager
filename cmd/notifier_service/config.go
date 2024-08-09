package main

import (
	"github.com/rokn/notifications-manager/pkg/channels"
	"github.com/rokn/notifications-manager/pkg/config"
)

type notifierConfig struct {
	config.DefaultConfig
	ChannelsServerURL string               `env:"CHANNELS_SERVER_URL" envDefault:"localhost:8080"`
	RabbitURI         string               `env:"RABBIT_URI" envDefault:"amqp://guest:guest@localhost:5672/"`
	ChannelType       channels.ChannelType `env:"CHANNEL_TYPE" envDefault:"email"`
}
