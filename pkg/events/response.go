package events

const (
	ResponseEventType = "botless.slack.response"
)

// Message, start simple.
type Message struct {
	Channel string `json:"channel,omitempty"`
	Text    string `json:"text,omitempty"`
}
