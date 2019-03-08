package events

const (
	MessageEventType  = "botless.slack.message"
	ResponseEventType = "botless.slack.response"
	WelcomeEventType  = "botless.slack.welcome"
)

// Message, start simple.
type Message struct {
	Channel string `json:"channel,omitempty"`
	Text    string `json:"text,omitempty"`
}
