package main

import (
	"log"
	"os"

	"github.com/botless/slack/pkg/slack"
	"github.com/kelseyhightower/envconfig"
)

// https://api.slack.com/slack-apps
// https://api.slack.com/internal-integrations
type envConfig struct {
	// Port is server port to be listened.
	Port int `envconfig:"BOT_PORT" default:"8080"`

	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// VerificationToken is used to validate interactive messages from slack.
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID"` // required:"true"`

	// ChannelID is slack channel ID where bot is working.
	// Bot responses to the mention in this channel.
	ChannelID string `envconfig:"CHANNEL_ID"` // required:"true"`

	// Sink is the consumer of cloud events from the bot.
	Sink string `envconfig:"SINK" required:"true"`
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	s := slack.New(env.BotToken, env.ChannelID, env.Sink, env.Port)

	if err := <-s.Err; err != nil {
		log.Printf("[ERROR] slack returned an error: %s", err)
		return 1
	}
	return 0
}

/*

if strings.HasPrefix(ev.Text, "echo ") {
				s.rtm.SendMessage(s.rtm.NewOutgoingMessage(
					strings.Replace(ev.Text, "echo ", "", 1),
					s.Channel))
			}

*/
