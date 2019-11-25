package bot

import (
	"fmt"
	"github.com/botless/events/pkg/events"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
)

func (s *Bot) cloudEventReceiver(event cloudevents.Event) {
	switch event.Type() {
	case "botless.bot.response":
		// don't block the cloudevents client.
		go func() {
			s.doResponse(event)
		}()
	}
}

func (s *Bot) doResponse(event cloudevents.Event) {
	resp := events.Message{}
	if err := event.DataAs(&resp); err != nil {
		s.Err <- fmt.Errorf("failed to get data from cloudevent %s", event.String())
	}
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(resp.Text, resp.Channel))
}
