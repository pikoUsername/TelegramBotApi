package main

import (
	"fmt"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/fsm/storage"
)

// main function entry function for whole program
func main() {
	bot, err := tgp.NewBot("<token>", "HTML", nil)

	// check out for error
	if err != nil {
		panic(err)
	}

	dp := tgp.NewDispatcher(bot, storage.NewMemoryStorage())

	// register a your callback
	// this callback will be called on every message
	// because handler havenot got any filters
	dp.MessageHandler.HandlerFunc(func(ctx *tgp.Context) {
		// returning message text to same chat
		_, err := ctx.Reply(tgp.NewReplyMessage(ctx.Message.Text))
		if err != nil {
			// you can use a more complex logging systems
			// it s just example
			fmt.Println(err)
		}
	})

	// if your bot has a payment or something important, then put skip_updates on false
	dp.RunPolling(tgp.NewPollingConfig(true))
}
