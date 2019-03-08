package main

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	clienthttp "github.com/cloudevents/sdk-go/pkg/cloudevents/client/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"log"
)

const (
	sourceURL = "https://github.com/knative/pkg#cloudevents-example"
	eventType = "botless.cloudevents.demo"
)

type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func defaultEventFields(event cloudevents.Event) cloudevents.Event {
	// get the context
	var ctx cloudevents.EventContextV02
	if event.Context != nil {
		ctx = event.Context.AsV02()
	} else {
		ctx = cloudevents.EventContextV02{}
	}

	// set the defaults
	ctx.Source = *types.ParseURLRef(sourceURL)
	ctx.Type = eventType
	// set it back
	event.Context = ctx
	return event
}

func main() {
	c, err := clienthttp.New(
		http.WithTarget("http://localhost:8080"),
		http.WithBinaryEncoding(),
		client.WithEventDefaulter(defaultEventFields),
		client.WithTimeNow(),
		client.WithUUIDs(),
	)
	if err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	for i := 0; i < 10; i++ {
		data := &Example{
			Message:  "hello, world!",
			Sequence: i,
		}

		if _, err := c.Send(context.TODO(), cloudevents.Event{Data: data}); err != nil {
			fmt.Printf("failed to send cloudevent: %v\n", err)
		}
	}
}
