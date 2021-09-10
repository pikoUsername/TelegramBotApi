package tgp_test

import (
	"fmt"
	"testing"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/filters"
	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

func TestRegister(t *testing.T) {
	dp, err := GetDispatcher(false)
	FailIfErr(t, err)
	// Simple echo handler
	dp.MessageHandler.Register(func(m *objects.Message) {
		bot := dp.Bot
		msg, err := bot.Send(&tgp.SendMessageConfig{
			ChatID: int64(m.From.ID),
			Text:   m.Text,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.Text)
	}, filters.CommandStart())
}

func TestMiddlwareRegister(t *testing.T) {
	dp, err := GetDispatcher(false)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	// this middleware will be a pre-process middleware
	// func(u *objects.Update) error/bool {...} will be a process middleware
	// and last middleware type is post process, maybe will be in this type
	// func(u objects.Update) {...}
	dp.MessageHandler.RegisterMiddleware(func(u *objects.Update) {
		// You can write any stuff you want to
		// FOr example simple ACL, or maybe other
		dp.Storage.SetData(
			u.Message.Chat.ID,
			u.Message.From.ID,
			storage.PackType{"AAAAAAAAAAA": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAa"},
		)
	})
}
