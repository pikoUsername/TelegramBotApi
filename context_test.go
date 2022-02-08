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

func GetContext(t *testing.T) *Context {
	dp, err := GetDispatcher(false)
	if err != nil {
		t.Fatal(err)
	}
	return dp.Context(fakeUpd)
}

func TestContextNext(t *testing.T) {
	ctx := GetContext(t)

	var x = 0

	h := NewHandlerType(func(c *Context) { x += 1 })
	ctx.handlers = append(ctx.handlers, h)

	ctx.Next()
	if x != 1 {
		t.Fatal("Handler didnt called")
	}
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
