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

type FileableConf interface {
	Configurable
}

// For CopyMessage method config
// https://core.telegram.org/bots/api#copymessage
type CopyMessageConfig struct {
	ChatID                int64 // required
	FromChatID            int64 // required
	MessageID             int64 // required
	Caption               string
	CaptionEntities       []*objects.MessageEntity
	DisableNotifications  bool
	ReplyToMessageId      int64
	AllowSendingWithReply bool

	// Note: Interface here is simple need
	// type: Union[objects.InlineKeyboardMarkup, ReplyKeyboardMarkup]
	ReplyMarkup interface{}
}

func (cmc *CopyMessageConfig) values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(cmc.ChatID, 10))
	v.Add("from_chat_id", strconv.FormatInt(cmc.ChatID, 10))
	v.Add("message_id", strconv.FormatInt(cmc.MessageID, 10))
	// TODO: Make Optional methods too...
	return v, nil
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
	return "sendMessage"
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

type GetUpdatesConfig struct {
	Offset         int
	Limit          int
	Timeout        int
	AllowedUpdates []string
}

func (guc *GetUpdatesConfig) values() (*url.Values, error) {
	v := &url.Values{}
	if guc.Offset != 0 {
		v.Add("offset", strconv.Itoa(guc.Offset))
	}
	v.Add("limit", strconv.Itoa(guc.Limit))
	v.Add("timeout", strconv.Itoa(guc.Timeout))

	return v, nil
}

func (guc *GetUpdatesConfig) method() string {
	return "getUpdates"
}
