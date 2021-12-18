package tgp

import (
	"testing"

	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

func GetDispatcher(check_token bool) (*Dispatcher, error) {
	var err error
	var b *Bot

	if check_token {
		b, err = NewBot(testToken, "HTML", nil)
	} else {
		b = &Bot{}
	}
	if err != nil {
		return &Dispatcher{}, err
	}
	return NewDispatcher(b, storage.NewMemoryStorage()), nil
}

func TestNewDispatcher(t *testing.T) {
	dp, _ := GetDispatcher(false)
	if dp == nil {
		t.Error("Oh no, Dispatcher didnt create, fix it")
		t.Fail()
	}
}

func TestProcessOneUpdate(t *testing.T) {
	dp, err := GetDispatcher(false)
	if err != nil {
		t.Fatal(err)
	}
	dp.ProcessOneUpdate(fakeUpd)
}

func TestProcessUpdates(t *testing.T) {
	dp, err := GetDispatcher(false)
	if err != nil {
		t.Fatal(err)
	}
	ch := make(chan *objects.Update)
	ch <- fakeUpd
	err = dp.ProcessUpdates(ch)
	if err != nil {
		t.Fatal(err)
	}
}

// go test -bench -benchmem

func BenchmarkProcessOneUpdate(b *testing.B) {
	dp, err := GetDispatcher(false)
	if err != nil {
		b.Error(err)
		b.Fail()
	}

	dp.MessageHandler.Register(func(ctx *Context) {})

	upd := &objects.Update{
		UpdateID: 100,
		Message:  &objects.Message{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dp.ProcessOneUpdate(upd)
	}
}
