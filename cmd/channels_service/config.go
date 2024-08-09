package main

import (
	"github.com/rokn/notifications-manager/pkg/config"
)

type channelsConfig struct {
	config.DefaultConfig
	ChannelsConfig string `env:"CHANNELS_CONFIG" envDefault:"./channels.yaml"`
}
