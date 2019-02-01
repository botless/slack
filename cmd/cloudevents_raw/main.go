package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/knative/pkg/cloudevents"
)

type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func handler(ctx context.Context, data json.RawMessage) {
	metadata := cloudevents.FromContext(ctx)
	log.Printf("[%s] %s %s: %q", metadata.EventTime.Format(time.RFC3339), metadata.ContentType, metadata.Source, string(data))
}

func main() {
	log.Print("ready and listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", cloudevents.Handler(handler)))
}
