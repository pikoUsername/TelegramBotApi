package bot

import "github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/objects"

// This file stores ALL method configs

// For CopyMessage method config
type CopyMessageConfig struct {
	ChatID                int64
	FromChatID            int64
	MessageID             int64
	Caption               string
	CaptionEntities       []*objects.MessageEntity
	DisableNotifications  bool
	ReplyToMessageId      int64
	AllowSendingWithReply bool

	// Note: Interface here is simple need
	// type: Union[objects.InlineKeyboardMarkup, ReplyKeyboardMarkup]
	ReplyMarkup interface{}
}
