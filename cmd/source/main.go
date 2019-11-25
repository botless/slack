package main

import (
	"context"
	bot2 "github.com/botless/slack/pkg/bot"
	"log"
)

func main() {
	ctx := context.Background()

	bot, err := bot2.New(ctx)
	if err != nil {
		log.Fatalf("failed to create bot: %s", err.Error())
	}
	if err := bot.Start(ctx); err != nil {
		log.Fatalf("failed to start cloudevent reciever: %s", err.Error())
	}
}
