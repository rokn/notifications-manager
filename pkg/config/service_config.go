package config

import (
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
	"log"
)

type ProfileType string

const (
	ProfileDevelopment ProfileType = "development"
	ProfileProduction  ProfileType = "production"
)

type ServiceConfig interface {
	GetPort() int
	GetProfile() ProfileType
}

type DefaultConfig struct {
	// Port to listen on
	Port int `env:"PORT" envDefault:"8080"`
	// Profile to run the service
	Profile ProfileType `env:"PROFILE" envDefault:"development"`
}

func (cfg *DefaultConfig) GetPort() int {
	return cfg.Port
}

func (cfg *DefaultConfig) GetProfile() ProfileType {
	return cfg.Profile
}

// InitConfigWithLogger Initialize a service configuration
func InitConfigWithLogger(config ServiceConfig) *zap.Logger {
	if err := env.Parse(config); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	return initLogger(config)
}

// Initialize the logger
func initLogger(cfg ServiceConfig) *zap.Logger {
	var logger *zap.Logger
	if cfg.GetProfile() == ProfileDevelopment {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	return logger
}
