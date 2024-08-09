package ingress

type NotificationDTO struct {
	// Name of the channels to send the notification to
	Channels []string `json:"channels" example:"email" binding:"min=1"`

	// Title of the notification
	Title string `json:"title" example:"Hello World"`

	// Body of the notification
	Message string `json:"message" example:"This is a test notification" binding:"required"`

	// Additional metadata
	Metadata map[string]string `json:"metadata"`
}

type NotificationResponse struct {
	// Message indicating the status of the notification
	Message string `json:"message"`
}

type ErrorResponse struct {
	// Error message
	Error string `json:"error"`
}

type ChannelsResponse struct {
	// List of channel names
	Channels []string `json:"channels"`
}
