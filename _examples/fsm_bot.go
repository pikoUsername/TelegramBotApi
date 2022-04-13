package main

import (
	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/filters"
	"github.com/pikoUsername/tgp/fsm"
	"github.com/pikoUsername/tgp/fsm/storage"
)

var (
	first_state  = fsm.NewState("HA")
	second_state = fsm.NewState("HO")
)

func main() {
	bot, err := tgp.NewBot("<token>", "HTML", nil)

	if err != nil {
		panic(err)
	}

	dp := tgp.NewDispatcher(bot, storage.NewMemoryStorage())

	dp.MessageHandler.HandlerFunc(func(ctx *tgp.Context) {
		ctx.Reply(tgp.NewSendMessage("Donald Duck is watching you."))
		ctx.SetState(first_state)
	})

	dp.MessageHandler.HandlerFunc(func(ctx *tgp.Context) {
		ctx.Reply(tgp.NewSendMessage("And big floppa too."))
		ctx.SetState(second_state)
	}.Filters(filters.StateFilter(first_state, dp.Storage))

	dp.RunPolling(tgp.NewPollingConfig(true))
}
