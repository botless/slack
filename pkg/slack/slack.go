package slack

import (
	"context"
	"fmt"
	"github.com/botless/slack/pkg/events"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/nlopes/slack"
	"log"
	"os"
	"strings"
)

type Slack struct {
	Channel string
	Err     chan error

	client *slack.Client
	rtm    *slack.RTM
	ce     client.Client

	domain string
}

const (
	source_template    = "https://%s.slack.com/messages/%s/" // domain, channel
	eventType_template = "botless.slack.*s"
)

type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func New(token, channel, target string, port int) *Slack {

	s := &Slack{
		Channel: channel,
		Err:     make(chan error, 1),
	}

	s.client = slack.New(
		token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	var err error
	if s.ce, err = client.NewHTTPClient(
		client.WithTarget(target),
		client.WithHTTPBinaryEncoding(),
		client.WithHTTPPort(port),
		client.WithTimeNow(),
		client.WithUUIDs(),
	); err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	// Use RTM:
	s.rtm = s.client.NewRTM()
	go s.rtm.ManageConnection()
	go s.manageRTM()

	if err = s.ce.StartReceiver(context.TODO(), s.cloudEventReceiver); err != nil {
		log.Fatalf("failed to start cloudevent reciever: %s", err.Error())
	}

	return s
}

func (s *Slack) manageRTM() {
	if team, err := s.client.GetTeamInfo(); err == nil {
		fmt.Printf("Slack Team: %+v", team)
	}

	for msg := range s.rtm.IncomingEvents {
		fmt.Println("Event Received: ", msg.Type)

		eventType := strings.ToLower(fmt.Sprintf(eventType_template, msg.Type))

		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			source := types.ParseURLRef(fmt.Sprintf(source_template, s.domain, ev.Channel))

			if err := s.ce.Send(context.TODO(), cloudevents.Event{
				Context: cloudevents.EventContextV02{
					Type:   eventType,
					Source: *source,
				},
				Data: ev,
			}); err != nil {
				fmt.Printf("failed to send cloudevent: %v\n", err)
			}

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())
			s.Err <- fmt.Errorf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			s.Err <- fmt.Errorf("invalid credentials")

		default:
			// Ignore other events..
			fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func (s *Slack) cloudEventReceiver(event cloudevents.Event) {
	resp := events.Message{}
	if err := event.DataAs(&resp); err != nil {
		s.Err <- fmt.Errorf("failed to get data from cloudevent %s", event.String())
	}
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(resp.Text, resp.Channel))
}
