package cloudevents

import (
	"github.com/knative/pkg/cloudevents"
	"log"
)

type Client struct {
	Client *cloudevents.Client
	send   chan interface{}
	done   chan bool
}

const (
	MAX_SEND_CHANNEL = 10
)

func NewClient(eventType, source, target string) *Client {
	c := &Client{}
	c.Client = cloudevents.NewClient(target, cloudevents.Builder{
		Source:    source,
		EventType: eventType,
	})
	return c
}

// Channel returns a channel that can be used to invoke Client.Send via a chan.
// This method has a side effect of another thread to monitor the send channel.
// Call client.Close() to shutdown the monitor thread.
// Experimental, error handling is not fully developed.
func (c *Client) Channel() chan<- interface{} {
	if c.send == nil {
		c.done = make(chan bool)
		c.send = make(chan interface{}, MAX_SEND_CHANNEL)
		go c.monitorSend()
	}
	return c.send
}

// Close ends the channel monitor produced by calling Channel()
// Experimental, error handling is not fully developed.
func (c *Client) Close() {
	if c.send == nil {
		return
	}
	c.done <- true
	close(c.send)
	c.send = nil
}

// monitorSend is the thread that will watch the send channel and call
// client.Sent() with the provided data struct. It will exit if the send channel
// closes or something is received on done channel.
func (c *Client) monitorSend() {
	for {
		select {
		case data, ok := <-c.send:
			if ok == false {
				break
			}
			if err := c.Client.Send(data); err != nil {
				log.Printf("error sending: %v", err)
			}
		case <-c.done:
			return
		}
	}
}
