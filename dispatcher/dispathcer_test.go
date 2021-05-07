package dispatcher_test

import (
	"testing"

	"github.com/pikoUsername/tgp/bot"
	"github.com/pikoUsername/tgp/dispatcher"
)

const (
	TestToken = "1780004238:AAENHJU4i9PaSIkgNjw-P2OvcQrtrO96JB4"
)

func GetDispatcher(t *testing.T) (*dispatcher.Dispatcher, error) {
	b, err := bot.NewBot(TestToken, true, "HTML")
	if err != nil {
		t.Error(err)
	}
	return &dispatcher.Dispatcher{Bot: b}, nil
}

func TestNewDispatcher(t *testing.T) {
	dp, err := GetDispatcher(t)
	if err != nil {
		t.Error(err)
	}
	if dp == nil {
		t.Error("Oh no, Dispatcher didnt create, fix it")
	}
}
