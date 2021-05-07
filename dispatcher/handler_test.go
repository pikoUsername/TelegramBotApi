package dispatcher_test

import (
	"fmt"
	"testing"

	"github.com/pikoUsername/tgp/configs"
	"github.com/pikoUsername/tgp/dispatcher"
	"github.com/pikoUsername/tgp/objects"
)

func TestRegister(t *testing.T) {
	dp, err := GetDispatcher(t)
	if err != nil {
		t.Error(err)
	}
	// Simple echo handler
	dp.MessageHandler.Register(func(upd objects.Update) {
		bot := dp.Bot
		msg, err := bot.Send(&configs.SendMessageConfig{
			ChatID: int64(upd.Message.From.ID),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.Text)
	})
}

func TestMiddlwareRegister(t *testing.T) {
	dp, err := GetDispatcher(t)
	if err != nil {
		t.Error(err)
	}
	dp.MessageHandler.RegisterMiddleware(func(u *objects.Update, hf dispatcher.HandlerFunc) {
		// You can write any stuff you want to
		// FOr example simple ACL, or maybe other
		if u.Message.From.FirstName == "Aleksei" {
			hf(*u)
		}
	})
}
