package main

import (
	"log"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/fsm/storage"
)

func main() {
	bot, err := tgp.NewBot("<token>", "HTML")
	var admin_id int64 = 10010101

	// check out for error
	if err != nil {
		panic(err)
	}
	dp := tgp.NewDispatcher(bot, storage.NewMemoryStorage())
	dp.MessageHandler.Register(func(c *tgp.Context) {
		if c.User.ID != admin_id {
			c.Next()
			return
		}
		c.Reply(tgp.NewSendMessage("Welcome, to the city 17."))
	})
	dp.MessageHandler.Register(func(c *tgp.Context) {
		c.Reply(tgp.NewSendMessage("Welcome to the Moscow"))
	})
	log.Fatal(dp.RunPolling(tgp.NewPollingConfig(true)))
}
