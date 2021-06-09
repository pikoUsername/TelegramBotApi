package dispatcher_test

import (
	"os"
	"testing"

	"github.com/pikoUsername/tgp/bot"
	"github.com/pikoUsername/tgp/dispatcher"
)

var (
	TestToken = os.Getenv("TEST_TOKEN")
)

func GetDispatcher(t *testing.T) *dispatcher.Dispatcher {
	b, err := bot.NewBot(TestToken, true, "HTML")
	if err != nil {
		t.Error(err)
	}
	return &dispatcher.Dispatcher{Bot: b}
}

func TestNewDispatcher(t *testing.T) {
	dp := GetDispatcher(t)
	if dp == nil {
		t.Error("Oh no, Dispatcher didnt create, fix it")
	}
}
