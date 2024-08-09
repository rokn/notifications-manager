package channels

type ChannelType string

const (
	Slack ChannelType = "slack"
	Email ChannelType = "email"
	SMS   ChannelType = "sms"
)

func (n ChannelType) In(s []ChannelType) bool {
	for _, v := range s {
		if v == n {
			return true
		}
	}
	return false
}

// ChannelModel data for a notification channel
type ChannelModel struct {
	// Name of the channel
	Name string
	// Type of the channel
	Type ChannelType
	// Configuration for the channel
	Configuration map[string]string
}
