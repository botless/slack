package main

import (
	"github.com/botless/slack/pkg/cloudevents"
	"log"
)

type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func main() {
	c := cloudevents.NewClient(
		"dev.knative.cloudevent.example",
		"https://github.com/knative/pkg#cloudevents-example",
		"http://localhost:8080",
	)

	for i := 0; i < 10; i++ {
		data := &Example{
			Message:  "hello, world!",
			Sequence: i,
		}
		if err := c.Send(data); err != nil {
			log.Printf("error sending: %v", err)
		}
	}
}
