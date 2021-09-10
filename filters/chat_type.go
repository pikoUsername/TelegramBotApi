package filters

import "github.com/pikoUsername/tgp/objects"

type ChatTypeFilter struct {
	Ignore string
}

func (ct *ChatTypeFilter) Check(u *objects.Update) bool {
	var chat *objects.Chat

	if u.EditedMessage != nil {
		chat = u.EditedMessage.Chat
	} else if u.ChannelPost != nil {
		chat = u.ChannelPost.Chat
	} else if u.Message != nil {
		chat = u.Message.Chat
	} else {
		chat = &objects.Chat{}
	}
	return chat.Type == ct.Ignore
}

func ChatType(ig string) *ChatTypeFilter {
	return &ChatTypeFilter{
		Ignore: ig,
	}
}
