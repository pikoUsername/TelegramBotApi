package fsm

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

	return group + ":" + s.State
}

// For NewState("...").Group("...")
func (s *State) Group(group string) *State {
	s.GroupState = group
	return s
}

// NewState init function
func NewState(state string) *State {
	return &State{
		State: state,
	}
}

// Uses for any purpose
var (
	DefaultState = NewState("")
	AnyState     = NewState("*")
)
