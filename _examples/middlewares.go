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
	dp.MessageHandler.HandlerFunc(func(c *tgp.Context) {
		if c.User.ID != admin_id {
			c.Next()
			return
		}
		c.Reply(tgp.NewSendMessage("Welcome, to the city 17."))
	})
	dp.MessageHandler.HandlerFunc(func(c *tgp.Context) {
		c.Reply(tgp.NewSendMessage("Welcome to the Moscow"))
	})
	// OR second variant
	// for single handlers(count it as simplyfied filters)
	dp.MessageHandler.HandlerFunc(func(ctx *tgp.Context) {
		ctx.Reply(tgp.NewReplyMessage("Hello admin"))
	}).Use(func(hand tgp.HandlerFunc) tgp.HandlerFunc {
		return func(ctx *tgp.Context) {
			if ctx.Message.From.ID == admin_id {
				hand(ctx)
				return
			} else {
				ctx.Reply(tgp.NewReplyMessage("You are not admin"))
			}
		}
	})
	// OR third variant
	// for all handlers that in specified handler group
	dp.MessageHandler.Use(func(next tgp.HandlerFunc) tgp.HandlerFunc {
		return func(ctx *tgp.Context) {
			if ctx.Message.From.ID == admin_id {
				next(ctx)
				ctx.Reply(tgp.NewReplyMessage("Hello admin."))
				return
			} else {
				return
			}
		}
	})
	log.Fatal(dp.RunPolling(tgp.NewPollingConfig(true)))
}
