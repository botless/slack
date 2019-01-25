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

	data := &Example{
		Message: "hello, world!",
	}

	for i := 0; i < 10; i++ {
		data.Sequence = i
		if !c.Send(data) {
			log.Printf("error sending: %v", c.SendError)
		}
	}
}
