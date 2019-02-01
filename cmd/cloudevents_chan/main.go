package main

import (
	"github.com/botless/slack/pkg/cloudevents"
	"time"
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
	ch := c.Channel()

	data := &Example{
		Message: "hello, world!",
	}

	for i := 0; i < 10; i++ {
		data.Sequence = i
		ch <- *data
	}

	time.Sleep(1 * time.Second)
	c.Close()
}
