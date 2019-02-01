package events

const (
	ResponseEventType = "botless.slack.response"
)

// Response, start simple.
type Response struct {
	Channel string `json:"channel,omitempty"`
	Text    string `json:"text,omitempty"`
}
