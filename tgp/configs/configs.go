package configs

import (
	"net/url"
	"strconv"

	"github.com/pikoUsername/tgp/tgp/objects"
)

// This file stores ALL method configs

// Configurable is interface for using by method
type Configurable interface {
	values() (*url.Values, error)
	method() string
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

// SendMessageConfig respresnests method,
// and fields of sendMessage method of telegram
// https://core.telegram.org/bots/api#sendmessage
type SendMessageConfig struct {
	// Required Field
	ChatID int

	// It s too, Telegram excepts
	Text                  string
	ParseMode             string
	Entities              []*objects.MessageEntity
	DisableWebPagePreview bool
}

// values ...
func (smc *SendMessageConfig) values() (*url.Values, error) {
	result := &url.Values{}
	result.Add("chat_id", strconv.Itoa(smc.ChatID))
	result.Add("text", smc.Text)
	if smc.ParseMode != "" {
		result.Add("parse_mode", smc.ParseMode)
	}
	result.Add("disable_web_page_preview", strconv.FormatBool(smc.DisableWebPagePreview))
	return result, nil
}

// method ...
func (smc *SendMessageConfig) method() string {
	return "SendMessage"
}

// SendPhotoConfig ...
type SendPhotoConfig struct {
}

func (spc *SendPhotoConfig) values() (*url.Values, error) {
	return nil, nil
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
