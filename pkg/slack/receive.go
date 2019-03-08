package slack

import (
	"fmt"
	"github.com/botless/slack/pkg/events"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
)

func (s *Slack) cloudEventReceiver(event cloudevents.Event) {
	switch event.Context.GetType() {
	case events.ResponseEventType:
		s.doResponse(event)
	}
}

func (s *Slack) doResponse(event cloudevents.Event) {
	// don't block the cloudevents client.
	go func() {
		resp := events.Message{}
		if err := event.DataAs(&resp); err != nil {
			s.Err <- fmt.Errorf("failed to get data from cloudevent %s", event.String())
		}
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(resp.Text, resp.Channel))
	}()
}
