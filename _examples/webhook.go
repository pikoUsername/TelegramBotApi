package main

import (
	"fmt"
	"log"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

func failIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	bot, err := tgp.NewBot("<your token>", "HTML", nil)

	failIfErr(err)
	dp := tgp.NewDispatcher(bot, storage.NewMemoryStorage())

	dp.MessageHandler.Register(func(ctx *tgp.Context) {
		_, err := ctx.Reply(tgp.NewSendMessage(ctx.Message.Text))
		if err != nil {
			fmt.Println(err)
		}
	})

	conf := tgp.NewWebhookConfig("/<some secret key, token, randomly generated string and etc.>/", "0.0.0.0:8443")

	conf.Certificate, err = objects.NewInputFile("./localhost.crt", "certificate")
	failIfErr(err)
	conf.KeyFile, err = objects.NewInputFile("./localhost.key", "<none>")
	failIfErr(err)

	// URI must be secret, and known only by telegram, and you.
	log.Fatal(dp.RunWebhook(conf))
}
