package tgp

import (
	"testing"
)

func TestContextHandlerAbility(t *testing.T) {
	dp, err := GetDispatcher(false)
	if err != nil {
		t.Fatal(err)
	}
	ctx := dp.Context(fakeUpd)
	t.Error(ctx)
}
