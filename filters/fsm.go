package filters

import (
	"reflect"

	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

// StateFilter uses for filter a state
type StateFilter struct {
	Storage storage.Storage
	State   interface{}
}

func (sf *StateFilter) GetState(u *objects.Update) string {
	var cid, uid int64

	if u.Message != nil {
		cid = u.Message.Chat.ID
		uid = u.Message.From.ID
	} else if u.EditedMessage != nil {
		cid = u.EditedMessage.Chat.ID
		uid = u.EditedMessage.From.ID
	} else if u.Message.PinnedMessage != nil {
		cid = u.Message.PinnedMessage.Chat.ID
		uid = u.Message.PinnedMessage.From.ID
	} else if u.Message.ReplyToMessage != nil {
		cid = u.Message.ReplyToMessage.Chat.ID
		uid = u.Message.ReplyToMessage.From.ID
	}

	state, err := sf.Storage.GetState(cid, uid)
	if err != nil {
		return ""
	}
	return state
}

func (sf *StateFilter) checkState(state string) bool {
	if reflect.TypeOf(sf.State).Comparable() && sf.State == state {
		return true
	}

	switch st := sf.State.(type) {
	case string:
		return st == state
	case []string:
		for _, s := range st {
			if s == state {
				return true
			}
		}
	}

	return false
}

func (sf *StateFilter) Check(u *objects.Update) bool {
	state := sf.GetState(u)
	return sf.checkState(state) || state == "*"
}

func NewStateFilter(state struct{}) *StateFilter {
	return &StateFilter{
		State: state,
	}
}
