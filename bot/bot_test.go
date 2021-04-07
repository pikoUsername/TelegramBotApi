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
	b, err := bot.NewBot(TestToken, false, ParseMode)
	if err != nil {
		return b, err
	}
	return b, nil
}

func TestCheckToken(t *testing.T) {
	b, err := bot.NewBot("bla:bla:bla", true, "HTML")
	if err != nil && b == nil {
		t.Error(err)
		t.Fail()
	}
}

func TestGetUpdates(t *testing.T) {
	b, _ := bot.NewBot(TestToken, false, "HTML")
	_, err := b.GetUpdates(&configs.GetUpdatesConfig{})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestParseMode(t *testing.T) {
	b, _ := getBot(t)
	_, err := utils.Link("https://www.google.com", "lol")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	m := &configs.SendMessageConfig{
		ChatID: TestChatID,
		Text:   "aaaa",
	}
	_, err = b.SendMessageable(m)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
