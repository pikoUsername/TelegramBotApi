package fsm

import (
	"fmt"
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
