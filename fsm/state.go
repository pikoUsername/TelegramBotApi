package fsm

import (
	"fmt"

	"github.com/pikoUsername/tgp/fsm/storage"
)

// State ...
type State struct {
	State      string
	GroupState string
}

// GetFullState just creates string, with {GroupState}:{StateName} template
func (s *State) GetFullState() string {
	var group string

	if s.State == "*" || s.State == "" {
		return s.State
	}

	if s.GroupState != "" {
		group = s.GroupState
	} else {
		group = "@"
	}

	return fmt.Sprintf("%s:%s", group, s.State)
}

func (s *State) Set(storage storage.Storage) {

}

// NewState init function
func NewState(state string, groupstate string) *State {
	return &State{
		State:      state,
		GroupState: groupstate,
	}
}

// Uses for any purpose
var (
	DefaultState = NewState("", "")
	AnyState     = NewState("*", "")
)
