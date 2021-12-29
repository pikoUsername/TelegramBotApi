package filters

import (
	"github.com/pikoUsername/tgp/fsm"
	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

// StateFilter uses for filter a state
type FSMStateFilter struct {
	Storage storage.Storage
	State   *fsm.State
}

func (sf *FSMStateFilter) GetState(u *objects.Update) string {
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

func (sf *FSMStateFilter) checkState(state string) bool {
	return sf.State.GetFullState() == state
}

func (sf *FSMStateFilter) Check(u *objects.Update) bool {
	state := sf.GetState(u)
	h := sf.checkState(state) || state == "*"
	return h
}

// StateFilter state argument types - State
func StateFilter(state *fsm.State, storage storage.Storage) *FSMStateFilter {
	return &FSMStateFilter{
		State:   state,
		Storage: storage,
	}
}
