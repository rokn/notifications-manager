package queue

import "github.com/rokn/notifications-manager/pkg/channels"

type PublishRequest struct {
	// Channels to publish the notification to grouped by type (key: type, value: list of channels)
	Channels map[channels.ChannelType][]string

	// Message to be published
	Message NotificationMessage
}

type NotificationMessage struct {
	// Title of the notification
	Title string `json:"title"`

	// Body of the notification
	Message string `json:"message"`

	// Additional metadata
	Metadata map[string]string `json:"metadata"`
}
