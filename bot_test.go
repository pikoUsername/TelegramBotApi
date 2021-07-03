package tgp_test

import (
	"fmt"
	"testing"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/objects"
	"github.com/pikoUsername/tgp/utils"
)

const (
	ParseMode  = "HTML"
	TestChatID = -534916942
)

func getBot(t *testing.T) *tgp.Bot {
	b, err := tgp.NewBot(TestToken, true, ParseMode)
	if err != nil {
		t.Error(err)
	}
	b.Debug = true
	return b
}

func TestCheckToken(t *testing.T) {
	b, err := tgp.NewBot("bla:bla", true, "HTML")
	if err != nil && b == nil {
		t.Error(err)
	}
}

func TestGetUpdates(t *testing.T) {
	b, err := tgp.NewBot(TestToken, false, "HTML")
	if err != nil {
		t.Error(err)
	}
	_, err = b.GetUpdates(&tgp.GetUpdatesConfig{})
	if err != nil {
		t.Error(err)
	}
}

func TestParseMode(t *testing.T) {
	b := getBot(t)
	line, err := utils.Link("https://www.google.com", "lol")
	if err != nil {
		t.Error(err)
	}
	m := &tgp.SendMessageConfig{
		ChatID: TestChatID,
		Text:   line,
	}
	_, err = b.SendMessageable(m)
	if err != nil {
		t.Error(err)
	}
}

func TestSetWebhook(t *testing.T) {
	b := getBot(t)
	resp, err := b.SetWebhook(tgp.NewSetWebhook("<URL>"))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
}

func TestSetCommands(t *testing.T) {
	// NOT OK, FAILS
	b := getBot(t)
	cmd := &objects.BotCommand{Command: "31321", Description: "ALLOO"}
	ok, err := b.SetMyCommands(tgp.NewSetMyCommands(cmd))
	if err != nil {
		t.Error(err, ok)
	}
	cmds, err := b.GetMyCommands(&tgp.GetMyCommandsConfig{})
	if err != nil {
		t.Error(err, cmds)
	}
	t.Log(cmds)
	for _, c := range cmds {
		if c.Command == cmd.Command && c.Description == cmd.Description {
			t.Skip("Original: ", cmd, "From tg: ", c)
			return
		}
	}
	t.Error("Command which getted from telegram, is not same as original, Original: ", cmd)
}
