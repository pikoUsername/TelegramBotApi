package bot_test

import (
	"testing"

	"github.com/pikoUsername/tgp/bot"
	"github.com/pikoUsername/tgp/configs"
	"github.com/pikoUsername/tgp/utils"
)

const (
	TestToken  = "1780004238:AAGFsgj2pxzXWoUqn25YohCEb1ENKIQOr1Q" // PolshaStrong_test_bot, yeah
	ParseMode  = "HTML"
	TestChatID = -534916942
)

func getBot(t *testing.T) (*bot.Bot, error) {
	b, err := bot.NewBot(TestToken, true, ParseMode)
	if err != nil {
		return b, err
	}
	b.Debug = true
	return b, nil
}

func TestCheckToken(t *testing.T) {
	b, err := bot.NewBot("bla:bla", true, "HTML")
	if err != nil && b == nil {
		t.Error(err)
	}
}

func TestGetUpdates(t *testing.T) {
	b, err := bot.NewBot(TestToken, false, "HTML")
	if err != nil {
		t.Error(err)
	}
	_, err = b.GetUpdates(&configs.GetUpdatesConfig{})
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
	m := &configs.SendMessageConfig{
		ChatID: TestChatID,
		Text:   line,
	}
	_, err = b.SendMessageable(m)
	if err != nil {
		t.Error(err)
	}
}
