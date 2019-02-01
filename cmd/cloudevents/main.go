package main

import (
	"fmt"
	"github.com/knative/pkg/cloudevents"
)

type Example struct {
	Sequence int    `json:"id"`
	Message  string `json:"message"`
}

func main() {
	c := cloudevents.NewClient(
		"http://localhost:8080",
		cloudevents.Builder{
			Source:    "https://github.com/knative/pkg#cloudevents-example",
			EventType: "dev.knative.cloudevent.example",
		})

	for i := 0; i < 10; i++ {
		data := &Example{
			Message:  "hello, world!",
			Sequence: i,
		}

		if err := c.Send(data, cloudevents.V01EventContext{
			Extensions: map[string]interface{}{
				"example": "example_ext",
			},
		}); err != nil {
			fmt.Printf("failed to send cloudevent: %v\n", err)
		}
	}
}
