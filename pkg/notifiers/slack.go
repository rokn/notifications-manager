package notifiers

import (
	"fmt"
	"github.com/rokn/notifications-manager/pkg/channels"
	"github.com/rokn/notifications-manager/pkg/queue"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type slackNotifier struct {
	log     *zap.Logger
	clients map[string]*slack.Client
}

func NewSlackNotifier(logger *zap.Logger) Notifier {
	return &slackNotifier{
		log:     logger.With(zap.String("notifier", "slack")),
		clients: make(map[string]*slack.Client),
	}
}

func (s *slackNotifier) SendNotification(channel *channels.ChannelModel, notification queue.NotificationMessage) error {
	s.log.Info("sending Slack notification", zap.String("channel", channel.Name))
	return fmt.Errorf("not implemented")
	client, err := s.getSlackClient(channel)
	if err != nil {
		return err
	}
	if channel.Configuration["channel_id"] == "" {
		s.log.Error("missing channel_id in configuration")
		return fmt.Errorf("missing channel_id in configuration")
	}

	// Create the Slack attachment that we will send to the channel
	attachment := slack.Attachment{
		Title: notification.Title,
		Text:  notification.Message,
		Color: "#36a64f",
	}

	_, timestamp, err := client.PostMessage(
		channel.Configuration["channel_id"],
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		s.log.Error("failed to send message", zap.Error(err))
		return err
	}

	s.log.Info("message sent", zap.String("timestamp", timestamp))
	return nil
}

func (s *slackNotifier) SupportedChannelTypes() []channels.ChannelType {
	return []channels.ChannelType{channels.Slack}
}

func (s *slackNotifier) getSlackClient(channel *channels.ChannelModel) (*slack.Client, error) {
	if channel.Configuration["token"] == "" {
		s.log.Error("missing token in configuration")
		return nil, fmt.Errorf("missing token in configuration")
	}

	clientKey := fmt.Sprintf("%s-%s", channel.Name, channel.Configuration["token"])
	if client, found := s.clients[clientKey]; found {
		return client, nil
	}

	client := slack.New(channel.Configuration["token"])
	s.clients[channel.Name] = client

	return client, nil

}
