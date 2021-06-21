package tgp_test

import (
	"testing"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/utils"
)

const (
	ParseMode  = "HTML"
	TestChatID = -534916942
)

func getBot(t *testing.T) (*tgp.Bot, error) {
	b, err := tgp.NewBot(TestToken, true, ParseMode)
	if err != nil {
		return b, err
	}
	b.Debug = true
	return b, nil
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
	b, err := getBot(t)
	if err != nil {
		t.Error(err)
	}
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
