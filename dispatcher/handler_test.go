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
	dp.MessageHandler.RegisterMiddleware(func(u *objects.Update, ht dispatcher.HandlerType) {
		u.Message.Text = "жил был, один фанат исекая. А на завтра его сбила машина, КОНЕЦ(а он ждал, и ждал в бесонечной темноте, пока его личность не растворилась навсегда)!"
		ht(*u)
	})
}
