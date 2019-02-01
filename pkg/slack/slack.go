package slack

import (
	"context"
	"fmt"
	"github.com/botless/slack/pkg/events"
	"github.com/knative/pkg/cloudevents"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Slack struct {
	Channel string
	Err     chan error

	client *slack.Client
	rtm    *slack.RTM
	ce     *cloudevents.Client

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

func New(token, channel, target, port string) *Slack {

	s := &Slack{
		Channel: channel,
		Err:     make(chan error, 1),
	}

	s.client = slack.New(
		token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	s.ce = cloudevents.NewClient(target, cloudevents.Builder{
		EventTypeVersion: "v1alpha1",
	})

	// Use RTM:
	s.rtm = s.client.NewRTM()
	go s.rtm.ManageConnection()
	go s.manageRTM()

	// CloudEvents incoming request // TODO: move this out.
	go s.manageIncomingCloudEvents(port)

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
			source := fmt.Sprintf(source_template, s.domain, ev.Channel)

			if err := s.ce.Send(ev, cloudevents.V01EventContext{
				EventType: eventType,
				Source:    source,
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

func (s *Slack) manageIncomingCloudEvents(port string) {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), cloudevents.Handler(s.handleIncomingCloudEvent)))
}

func (s *Slack) handleIncomingCloudEvent(ctx context.Context, resp *events.Response) {
	metadata := cloudevents.FromContext(ctx)
	log.Printf("[%s] %s %s: %q", metadata.EventTime.Format(time.RFC3339), metadata.ContentType, metadata.Source, resp.Text)

	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(resp.Text, resp.Channel))
}
