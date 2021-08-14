package tgp_test

import (
	"testing"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/objects"
)

func TestUnRegister(t *testing.T) {
	m := tgp.NewMiddlewareManager(nil)

	md := func(upd *objects.Update) {}

	m.Register(md)
	f, _ := m.Unregister(md)
	if f == nil {
		t.Error("tgp: Unregister is not working")
		t.Fail()
	}
}
