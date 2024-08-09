package channels

// Channel data for a notification channel
type Channel struct {
	// Name of the channel
	Name *string `yaml:"name" validate:"required"`
	// Type of the channel
	Type *ChannelType `yaml:"type" validate:"required"`
	// Configuration for the channel
	Configuration *map[string]string `yaml:"configuration"`
}

type Config struct {
	Channels []Channel `yaml:"channels" validate:"dive"`
}
