package main

import (
	"context"
	"github.com/botless/slack/pkg/events"
	"github.com/kelseyhightower/envconfig"
	"github.com/knative/pkg/cloudevents"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
	"strings"
)

type envConfig struct {
	// Port is server port to be listened.
	Port string `envconfig:"USER_PORT" default:"8080"`

	Target string `envconfig:"TARGET" required:"true"`
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

type Echo struct {
	ce *cloudevents.Client
}

func (e *Echo) handler(ctx context.Context, msg *slack.Message) {
	metadata := cloudevents.FromContext(ctx)
	_ = metadata

	log.Printf("Message: %s", msg.Text)

	if strings.HasPrefix(msg.Text, "echo ") {
		resp := events.Response{
			Channel: msg.Channel,
			Text:    strings.Replace(msg.Text, "echo ", "", 1),
		}
		if err := e.ce.Send(resp); err != nil {
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

	e := &Echo{
		ce: cloudevents.NewClient(env.Target, cloudevents.Builder{
			EventTypeVersion: "v1alpha1",
			EventType:        events.ResponseEventType,
			Source:           "slack.echo",
		}),
	}

	log.Print("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", cloudevents.Handler(e.handler)))
	return 0
}
