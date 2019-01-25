package slack

import (
	"fmt"
	"github.com/botless/slack/pkg/cloudevents"
	"github.com/nlopes/slack"
	"log"
	"os"
)

type Slack struct {
	Channel string
	Err     chan error

	client   *slack.Client
	rtm      *slack.RTM
	ceClient *cloudevents.Client
}

type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func New(token, channel, target string) *Slack {

	s := &Slack{
		Channel: channel,
		Err:     make(chan error, 1),
	}

	s.client = slack.New(
		token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	s.ceClient = cloudevents.NewClient(
		"dev.knative.cloudevent.example",
		"https://github.com/knative/pkg#cloudevents-example",
		target,
	)

	// Use RTM:
	s.rtm = s.client.NewRTM()
	go s.rtm.ManageConnection()
	go s.manageRTM()

	return s
}

func (s *Slack) manageRTM() {
	ce := s.ceClient.Channel()

	for msg := range s.rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			s.rtm.SendMessage(s.rtm.NewOutgoingMessage("Hello world", s.Channel))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

			me := slack.MessageEvent(*ev)

			ce <- Example{Message: fmt.Sprintf("%s: %s", me.User, me.Text)}

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
	s.ceClient.Done()
}
