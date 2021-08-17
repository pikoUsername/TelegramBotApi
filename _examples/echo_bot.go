package main

import (
	"fmt"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

// main function entry function for whole program
func main() {
	bot, err := tgp.NewBot("<token>", "HTML")

	// check out for error
	if err != nil {
		panic(err)
	}

	// recommended set syncronus argument to false
	dp := tgp.NewDispatcher(bot, storage.NewMemoryStorage(), false)

	// register a your callback
	// this callback will be called on every message
	// because handler havenot got any filters
	dp.MessageHandler.Register(func(bot *tgp.Bot, m *objects.Message) {
		// returning message text to same chat
		_, err := bot.Send(tgp.NewSendMessage(m.Chat.ID, m.Text))
		if err != nil {
			// you can use a more complex logging systems
			// it s just example
			fmt.Println(err)
		}
	})

	// if your bot has a payment or something important, then put skip_updates on false
	dp.StartPolling(tgp.NewStartPollingConf(true))
}
