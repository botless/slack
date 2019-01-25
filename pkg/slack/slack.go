package slack

import (
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"os"
)

type Slack struct {
	Channel string
	Err     chan error

	client *slack.Client
	rtm    *slack.RTM
}

func New(token, channel string) *Slack {

	s := &Slack{
		Channel: channel,
		Err:     make(chan error, 1),
	}

	s.client = slack.New(
		token,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	// Use RTM:
	s.rtm = s.client.NewRTM()
	go s.rtm.ManageConnection()

	return s
}

func (s *Slack) manageRTM() {
	for msg := range s.rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			s.rtm.SendMessage(s.rtm.NewOutgoingMessage("Hello world", s.Channel))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

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

/*

if env.ChannelID == "" {
		channels, err := rtm.GetChannels(true)
		if err != nil {
			log.Printf("channel error: %v", err)
		} else {
			log.Printf("channels: %+v", channels)
		}
		return 0
	}

*/
