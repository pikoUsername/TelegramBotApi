package tgp

import (
	"sync"
	"testing"

	"github.com/pikoUsername/tgp/fsm"
	"github.com/pikoUsername/tgp/fsm/storage"
)

var (
	testCtx = Context{
		Bot: nil,

		Update:  fakeUpd,
		Storage: nil,
		mu:      sync.Mutex{},
	}
)

func TestContextHandlerAbility(t *testing.T) {
	dp, err := GetDispatcher(false)
	if err != nil {
		t.Fatal(err)
	}
	ctx := dp.Context(fakeUpd)
	t.Error(ctx)
}

func TestResetState(t *testing.T) {
	testCtx.Storage = storage.NewMemoryStorage()
	testCtx.SetState(fsm.AnyState)
	s, err := testCtx.Storage.GetState(testCtx.Message.Chat.ID, testCtx.Message.From.ID)
	if err != nil {
		t.Fatal(err)
	}
	if s == "" {
		t.Fatal("No state")
	}
}
