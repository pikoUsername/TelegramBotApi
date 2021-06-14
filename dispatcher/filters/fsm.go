package filters

import (
	"github.com/pikoUsername/tgp/dispatcher/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

// StateFilter uses for filter a state
type StateFilter struct {
	Storage storage.Storage
	State   string
}

func (sf *StateFilter) GetState(u *objects.Update) string {
	var cid, uid int64

	if u.Message != nil {
		cid = u.Message.Chat.ID
		uid = u.Message.From.ID
	} else if u.EditedMessage != nil {
		cid = u.EditedMessage.Chat.ID
		uid = u.EditedMessage.From.ID
	} else if u.PinnedMessage != nil {
		cid = u.PinnedMessage.Chat.ID
		uid = u.PinnedMessage.From.ID
	} else if u.ReplyToMessage != nil {
		cid = u.ReplyToMessage.Chat.ID
		uid = u.ReplyToMessage.From.ID
	}

	return sf.Storage.GetState(cid, uid)
}

func (sf *StateFilter) Check(u *objects.Update) bool {
	state := sf.GetState(u)

	if state == "*" {
		return true
	}
	if state == sf.State {
		return true
	}

	return false
}
