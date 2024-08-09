package channels

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Service interface {
	GetChannel(name string) (*Channel, error)
	GetChannelNames() []string
}

type service struct {
	channels map[string]Channel
	log      *zap.Logger
}

func NewService(configPath string, logger *zap.Logger) Service {
	log := logger.With(zap.String("service", "channels"))

	// Load the configuration from yaml
	config := Config{}

	// Open the file
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal("failed to open config file", zap.Error(err))
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Read the file content
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("failed to read config file", zap.Error(err))
	}

	// Parse the YAML content
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("failed to parse config", zap.Error(err))
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		log.Fatal("invalid config", zap.Error(err))
	}

	// Create a map for quick lookup
	configs := make(map[string]Channel)
	for _, channel := range config.Channels {
		configs[*channel.Name] = channel
	}

	return &service{
		channels: configs,
		log:      log,
	}
}

func (svc *service) GetChannel(name string) (*Channel, error) {
	if channel, ok := svc.channels[name]; ok {
		return &channel, nil
	}
	return nil, errors.New("channel not found")
}

func (svc *service) GetChannelNames() []string {
	names := make([]string, 0, len(svc.channels))
	for name := range svc.channels {
		names = append(names, name)
	}
	return names
}
