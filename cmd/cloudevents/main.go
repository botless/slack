package main

import (
	"fmt"
	"github.com/botless/slack/pkg/cloudevents"
	ce "github.com/knative/pkg/cloudevents"
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

		if err := c.Client.Send(data, ce.V01EventContext{
			Extensions: map[string]interface{}{
				"example": "example_ext",
			},
		}); err != nil {
			fmt.Printf("failed to send cloudevent: %v\n", err)
		}

		if err := c.Client.Send(data); err != nil {
			log.Printf("error sending: %v", err)
		}
	}
}
