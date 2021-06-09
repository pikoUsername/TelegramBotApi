package dispatcher_test

import (
	"fmt"
	"testing"

	"github.com/pikoUsername/tgp/configs"
	"github.com/pikoUsername/tgp/dispatcher"
	"github.com/pikoUsername/tgp/objects"
)

func TestRegister(t *testing.T) {
	dp := GetDispatcher(t)
	// Simple echo handler
	dp.MessageHandler.Register(func(update *objects.Update) {
		bot := dp.Bot
		msg, err := bot.Send(&configs.SendMessageConfig{
			ChatID: int64(update.Message.From.ID),
			Text:   update.Message.Text,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.Text)
	}, nil)
}

func TestMiddlwareRegister(t *testing.T) {
	dp := GetDispatcher(t)
	dp.MessageHandler.RegisterMiddleware(func(u *objects.Update, hf dispatcher.HandlerType) {
		// You can write any stuff you want to
		// FOr example simple ACL, or maybe other
		if u.Message.From.FirstName == "Aleksei" {
			hf.Call(u)
		}
	})
}
