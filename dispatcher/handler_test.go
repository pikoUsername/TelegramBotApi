package dispatcher_test

import (
	"fmt"
	"testing"

	"github.com/pikoUsername/tgp/configs"
	"github.com/pikoUsername/tgp/dispatcher/filters"
	"github.com/pikoUsername/tgp/dispatcher/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

func TestRegister(t *testing.T) {
	dp := GetDispatcher(t)
	// Simple echo handler
	dp.MessageHandler.Register(func(m *objects.Message) {
		bot := dp.Bot
		msg, err := bot.Send(&configs.SendMessageConfig{
			ChatID: int64(m.From.ID),
			Text:   m.Text,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.Text)
	}, filters.NewCommandStart())
}

func TestMiddlwareRegister(t *testing.T) {
	dp := GetDispatcher(t)

	// this middleware will be a pre-process middleware
	// func(u *objects.Update) error/bool {...} will be a process middleware
	// and last middleware type is post process, maybe will be in this type
	// func(u objects.Update) {...}
	dp.MessageHandler.RegisterMiddleware(func(u *objects.Update) {
		// You can write any stuff you want to
		// FOr example simple ACL, or maybe other
		dp.Storage.SetData(
			u.Chat.ID, 
			u.From.ID,
			&storage.PackType{"AAAAAAAAAAA": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAa"}
		)
	})
}
