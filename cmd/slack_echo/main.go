package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/botless/slack/pkg/events"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	clienthttp "github.com/cloudevents/sdk-go/pkg/cloudevents/client/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	// Port is server port to be listened.
	Port int `envconfig:"USER_PORT" default:"8080"`

	// Target is the endpoint to receive cloudevents.
	Target string `envconfig:"TARGET" required:"true"`
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

type Echo struct {
	ce client.Client
}

func (e *Echo) receive(event cloudevents.Event) {
	msg := events.Message{}
	if err := event.DataAs(&msg); err != nil {
		log.Printf("failed to get data from cloudevent %s", event.String())
	}

	log.Printf("Message: %s", msg.Text)

	if strings.HasPrefix(msg.Text, "echo ") {
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				Type:   events.ResponseEventType,
				Source: *types.ParseURLRef("//botless/slack/echo"),
			},
			Data: events.Message{
				Channel: msg.Channel,
				Text:    strings.Replace(msg.Text, "echo ", "", 1),
			},
		}
		if _, err := e.ce.Send(context.TODO(), event); err != nil {
			log.Printf("failed to send cloudevent: %s\n", err)
		}
	}
}

func _main(args []string) int {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	ce, err := clienthttp.New(
		http.WithTarget(env.Target),
		http.WithPort(env.Port),
		http.WithBinaryEncoding(),
		client.WithTimeNow(),
		client.WithUUIDs(),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %s", err.Error())
	}
	e := &Echo{
		ce: ce,
	}

	ctx := context.Background()
	if err := e.ce.StartReceiver(ctx, e.receive); err != nil {
		log.Fatalf("Failed to create client: %s", err.Error())
	}
	log.Printf("listening on port %d", env.Port)
	<-ctx.Done()

	return 0
}
