package tgp_test

import (
	"testing"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

func GetDispatcher(check_token bool) (*tgp.Dispatcher, error) {
	var err error
	var b *tgp.Bot

	if check_token {
		b, err = tgp.NewBot(TestToken, "HTML")
	} else {
		b = &tgp.Bot{}
	}
	if err != nil {
		return &tgp.Dispatcher{}, err
	}
	return tgp.NewDispatcher(b, storage.NewMemoryStorage(), false), nil
}

func TestNewDispatcher(t *testing.T) {
	dp, _ := GetDispatcher(false)
	if dp == nil {
		t.Error("Oh no, Dispatcher didnt create, fix it")
		t.Fail()
	}
}

func BenchmarkProcessOneUpdate(b *testing.B) {
	dp, err := GetDispatcher(false)
	if err != nil {
		b.Error(err)
		b.Fail()
	}
	dp.MessageHandler.Register(func(m *objects.Message) {})
	dp.MessageHandler.RegisterMiddleware(func(u *objects.Update) {})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		upd := &objects.Update{
			UpdateID: i,
			Message:  &objects.Message{},
		}
		b.StartTimer()
		dp.ProcessOneUpdate(upd)
	}
}
