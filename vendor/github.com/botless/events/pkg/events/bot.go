package events

import (
	"fmt"
	"log"
	"strings"
)

const (
	bot_type_template = "botless.bot.%s" // type
)

var knownBotEvents = []string{
	"response", // payload: Message
	"command",  // payload: Command
}

type bot int

var Bot bot // export Bot
var _ = Bot

// Message, start simple.
type Message struct {
	Channel string `json:"channel,omitempty"`
	Text    string `json:"text,omitempty"`
}

type Command struct {
	Cmd     string `json:"cmd,omitempty"`
	Args    string `json:"args,omitempty"`
	Author  string `json:"author,omitempty"`
	Channel string `json:"channel,omitempty"`
}

// TODO: not sold on channel in Command. It should not be there. Nor should Author. I think these should be part of the source.
// Then ditto on Message.Channel

func (bot) Type(t ...string) string {
	if len(t) == 0 {
		return strings.ToLower(fmt.Sprintf(bot_type_template, "unknown"))
	}
	if !contains(knownBotEvents, t[0]) {
		log.Printf("[WARN] unknown bot event type: %q", t)
	}
	return strings.ToLower(fmt.Sprintf(bot_type_template, strings.Join(t, ".")))
}
