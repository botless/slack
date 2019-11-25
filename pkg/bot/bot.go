package bot

import (
	"context"
	"fmt"
	"log"
	"sync"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/kelseyhightower/envconfig"
	bs "github.com/mattmoor/bindings/pkg/slack"
	"github.com/nlopes/slack"
)

type Bot struct {
	// Sink is the consumer of cloud events from the bot.
	Sink string `envconfig:"SINK" required:"true"`
	// Port is server port to be listened.
	Port int `envconfig:"PORT" default:"8080"`

	Admin   string
	Bot     string
	Channel string
	Err     chan error

	client *slack.Client
	rtm    *slack.RTM
	ce     client.Client

	domain    string
	slackOnce sync.Once
}

const (
	channelKey = "channel"
	adminKey   = "admin"
	botKey     = "bot"
)

func New(ctx context.Context) (*Bot, error) {
	var bot Bot
	var err error

	if err := envconfig.Process("", &bot); err != nil {
		return nil, err
	}
	bot.Err = make(chan error, 1)

	if bot.Admin, err = bs.ReadKey(adminKey); err != nil {
		return nil, err
	}

	if bot.Bot, err = bs.ReadKey(botKey); err != nil {
		return nil, err
	}

	if bot.Channel, err = bs.ReadKey(channelKey); err != nil {
		return nil, err
	}

	if bot.client, err = bs.New(ctx); err != nil {
		return nil, err
	}

	t, err := cloudevents.NewHTTPTransport(
		http.WithTarget(bot.Sink),
		http.WithBinaryEncoding(),
		http.WithPort(bot.Port),
	)
	if err != nil {
		return nil, err
	}

	if bot.ce, err = cloudevents.NewClient(t,
		client.WithTimeNow(),
		client.WithUUIDs(),
	); err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	// Use RTM:
	bot.rtm = bot.client.NewRTM()
	go bot.rtm.ManageConnection()
	go bot.manageRTM(ctx)

	//ch, err := bot.client.GetChannelInfo(bot.Channel)
	//if err != nil {
	//	return nil, err
	//}

	user, err := bot.client.GetUserInfo(bot.Admin)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	fmt.Printf("ADMIN --> ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)

	return &bot, nil
}

// Start is a blocking call.
func (b *Bot) Start(ctx context.Context) error {
	if err := b.ce.StartReceiver(context.TODO(), b.cloudEventReceiver); err != nil {
		return err
	}
	return nil
}
