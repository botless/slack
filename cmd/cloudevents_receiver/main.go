package main

import (
	"context"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	clienthttp "github.com/cloudevents/sdk-go/pkg/cloudevents/client/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"log"
)

type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func receive(event cloudevents.Event) {
	var data map[string]interface{}
	if err := event.DataAs(&data); err != nil {
		log.Printf("failed to get data: %s", err.Error())
	}
	log.Printf("got %s, %+v", event.String(), data)
}

func main() {
	ctx := context.Background()

	c, err := clienthttp.New(http.WithPort(8080))
	if err != nil {
		log.Fatalf("Failed to create client: %s", err.Error())
	}

	if err := c.StartReceiver(ctx, receive); err != nil {
		log.Fatalf("Failed to start reveiver client: %s", err.Error())
	}
	log.Print("listening on port 8080")
	<-ctx.Done()
}
