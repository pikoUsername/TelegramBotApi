package main

import (
	"fmt"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/fsm/storage"
)

func main() {
	bot, err := tgp.NewBot("<token>", "HTML", nil)

	if err != nil {
		panic(err)
	}

	dp := tgp.NewDispatcher(bot, storage.NewMemoryStorage())

	dp.MessageHandler.Register(func(ctx *tgp.Context) {
		_, err := bot.Send(tgp.NewSendMessage(ctx.Message.Chat.ID, ctx.Message.Text))
		if err != nil {
			fmt.Println(err)
		}
	})

	dp.StartPolling(tgp.NewStartPollingConf(true))
}
