package fsm_test

import (
	"testing"

	"github.com/pikoUsername/tgp/fsm"
)

func TestGetFullState(t *testing.T) {
	// * is any state
	test_state := fsm.NewState("*")
	test_state1 := fsm.NewState("LOLLOLOLLOLLOLLOLLOLLOLLOLLOL").Group("kekek")

	// getting all states strings, huh
	fs := test_state.GetFullState()
	fs1 := test_state1.GetFullState()

	if fs != "*" {
		t.Error("* any state now is not", fs)
	}
	if fs1 != "kekek:LOLLOLOLLOLLOLLOLLOLLOLLOLLOL" {
		t.Error("Not correct string formation")
	}
}
