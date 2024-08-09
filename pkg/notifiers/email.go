package notifiers

import (
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/rokn/notifications-manager/pkg/channels"
	"github.com/rokn/notifications-manager/pkg/queue"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"strconv"
)

type emailNotifier struct {
	validate *validator.Validate
	log      *zap.Logger
}

type emailConfig struct {
	To       string `mapstructure:"to" validate:"required,email"`
	Host     string `mapstructure:"host" validate:"required"`
	Port     string `mapstructure:"port" validate:"required"`
	Username string `mapstructure:"username" validate:"required,email"`
	Password string `mapstructure:"password" validate:"required"`
}

func NewEmailNotifier(logger *zap.Logger) Notifier {
	return &emailNotifier{
		validate: validator.New(),
		log:      logger.With(zap.String("notifier", "email")),
	}
}

func (e *emailNotifier) SendNotification(channel *channels.ChannelModel, notification queue.NotificationMessage) error {
	e.log.Info("sending email notification", zap.String("channel", channel.Name), zap.Any("notification", notification))

	dialer, config, err := e.decodeConfig(channel)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Username)
	m.SetHeader("To", config.To)
	m.SetHeader("Subject", notification.Title)
	m.SetBody("text/html", notification.Message)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}

	e.log.Info("email notification sent", zap.String("channel", channel.Name))
	return nil
}

func (e *emailNotifier) SupportedChannelTypes() []channels.ChannelType {
	return []channels.ChannelType{channels.Email}
}

func (e *emailNotifier) decodeConfig(channel *channels.ChannelModel) (*gomail.Dialer, *emailConfig, error) {
	config := &emailConfig{}
	err := mapstructure.Decode(channel.Configuration, config)
	if err != nil {
		e.log.Error("failed to decode email configuration", zap.Error(err))
		return nil, nil, err
	}
	err = e.validate.Struct(config)
	if err != nil {
		e.log.Error("email configuration validation failed", zap.Error(err))
		return nil, nil, err
	}

	port, err := strconv.ParseInt(config.Port, 10, 64)
	if err != nil {
		e.log.Error("failed to parse port", zap.Error(err))
		return nil, nil, err
	}

	dialer := gomail.NewDialer(config.Host, int(port), config.Username, config.Password)
	return dialer, config, nil
}
