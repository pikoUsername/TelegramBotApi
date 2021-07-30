package tgp_test

import (
	"testing"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/fsm/storage"
)

func GetDispatcher(t *testing.T) *tgp.Dispatcher {
	b, err := tgp.NewBot(TestToken, true, "HTML", Timeout)
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

func TestOnStartup(t *testing.T) {
	dp := GetDispatcher(t)
	a := false
	dp.OnStartup(
		tgp.NewOnConf(
			func(dp *tgp.Dispatcher) {
				a = true
				dp.ShutDownDP()
			},
		),
	)
	if !a {
		t.Error("A variable, didnt change, fix it!")
	}
}
