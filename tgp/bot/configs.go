package bot

import (
	"net/url"

	"github.com/pikoUsername/tgp/tgp/objects"
)

// This file stores ALL method configs

// Configurable is interface for using by methods
type Configurable interface {
	values() *url.Values
}

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

type SendPhotoConfig struct {
}

type SendAudioConfig struct {
}

type SendDocumentConfig struct {
}

type SendVideoConfig struct {
}

type SendAnimation struct {
}

type SendVoiceConfig struct {
}

type SendVideoNameConfig struct {
}

type SendMediaGroupConfig struct {
}

type SendLocationConfig struct {
}

type LiveLocationConfig struct {
}
