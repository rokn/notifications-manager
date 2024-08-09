package main

import (
	"github.com/rokn/notifications-manager/pkg/config"
)

type ingressConfig struct {
	config.DefaultConfig
	ChannelsServerURL string `env:"CHANNELS_SERVER_URL" envDefault:"localhost:8080"`
	RabbitURI         string `env:"RABBIT_URI" envDefault:"amqp://guest:guest@localhost:5672/"`
}
