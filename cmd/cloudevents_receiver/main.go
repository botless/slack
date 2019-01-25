package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/knative/pkg/cloudevents"
)

type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func handler(ctx context.Context, data *Example) {
	metadata := cloudevents.FromContext(ctx)
	log.Printf("[%s] %s %s: %d,%q", metadata.EventTime.Format(time.RFC3339), metadata.ContentType, metadata.Source, data.Sequence, data.Message)
}

func main() {
	log.Print("ready and listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", cloudevents.Handler(handler)))
}
