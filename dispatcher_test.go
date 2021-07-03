package tgp_test

import (
	"os"
	"testing"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/fsm/storage"
)

var (
	TestToken = os.Getenv("TEST_TOKEN")
)

func GetDispatcher(t *testing.T) *tgp.Dispatcher {
	b, err := tgp.NewBot(TestToken, true, "HTML")
	if err != nil {
		t.Error(err)
	}
	return tgp.NewDispatcher(b, storage.NewMemoryStorage(), false)
}

func TestNewDispatcher(t *testing.T) {
	dp := GetDispatcher(t)
	if dp == nil {
		t.Error("Oh no, Dispatcher didnt create, fix it")
	}
}
