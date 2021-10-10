package tgp

import (
	"fmt"
	"testing"

	"github.com/pikoUsername/tgp/filters"
	"github.com/pikoUsername/tgp/objects"
)

var (
	// fake update
	fakeUpd = &objects.Update{
		Message: &objects.Message{
			MessageID: 1000,
			Chat: &objects.Chat{
				ID:        1000,
				FirstName: "LoL",
				Username:  "LoL",
			},
			From: &objects.User{
				ID:           1000,
				IsBot:        false,
				FirstName:    "KAK",
				LanguageCode: "ru",
				LastName:     "lol",
			},
			Text: "В",
		},
	}
)

func TestRegister(t *testing.T) {
	dp, err := GetDispatcher(false)
	failIfErr(t, err)
	// Simple echo handler
	dp.MessageHandler.Register(func(ctx *Context) {
		bot := dp.Bot
		msg, err := bot.Send(&SendMessageConfig{
			ChatID: int64(ctx.Message.From.ID),
			Text:   ctx.Message.Text,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.Text)
	}, filters.CommandStart())
}

func TestHandlerTrigger(t *testing.T) {
	dp, err := GetDispatcher(true)
	if err != nil {
		t.Fatal(err)
	}

	dp.MessageHandler.Register(func(ctx *Context) {
		fmt.Println("working!")
		ctx.Abort()
	})
	upd := &objects.Update{
		Message: &objects.Message{
			MessageID: 1000,
			Chat: &objects.Chat{
				ID:        1000,
				FirstName: "LoL",
				Username:  "LoL",
			},
			From: &objects.User{
				ID:           1000,
				IsBot:        false,
				FirstName:    "KAK",
				LanguageCode: "ru",
				LastName:     "lol",
			},
			Text: "В",
		},
	}
	dp.MessageHandler.Trigger(dp.Context(upd))
}
