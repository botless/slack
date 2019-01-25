package cloudevents

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	Builder
	Target string

	send chan interface{}
	done chan bool
}

const (
	MAX_SEND_CHANNEL = 10
)

func NewClient(eventType, source, target string) *Client {
	c := &Client{
		Builder: Builder{
			Source:    source,
			EventType: eventType,
		},
		Target: target,
	}
	return c
}

func (c *Client) RequestSend(data interface{}) (*http.Response, error) {
	req, err := c.Build(c.Target, data)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}

func (c *Client) Send(data interface{}) error {
	resp, err := c.RequestSend(data)
	if err != nil {
		return err
	}
	if Accepted(resp) {
		return nil
	}
	return fmt.Errorf("error sending cloudevent: %s", Status(resp))
}

func (c *Client) Channel() chan<- interface{} {
	if c.send == nil {
		c.done = make(chan bool)
		c.send = make(chan interface{}, MAX_SEND_CHANNEL)
		go c.monitorSend()
	}
	return c.send
}

func (c *Client) Done() {
	if c.send == nil {
		return
	}
	c.done <- true
	close(c.send)
	c.send = nil
}

func Accepted(resp *http.Response) bool {
	if resp.StatusCode == 204 {
		return true
	}
	return false
}

func Status(resp *http.Response) string {
	if Accepted(resp) {
		return "sent"
	}

	status := resp.Status
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Sprintf("Status[%s] error reading response body: %v", status, err)
	}

	return fmt.Sprintf("Status[%s] %s", status, body)
}

func (c *Client) monitorSend() {
	for {
		select {
		case data, ok := <-c.send:
			if ok == false {
				break
			}
			if err := c.Send(data); err != nil {
				log.Printf("error sending: %v", err)
			}
		case <-c.done:
			return
		}
	}
}
