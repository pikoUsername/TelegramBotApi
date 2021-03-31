package configs

import (
	"net/url"
	"strconv"

	"github.com/pikoUsername/tgp/tgp/objects"
)

// This file stores ALL method configs

// Configurable is interface for using by method
type Configurable interface {
	Values() (*url.Values, error)
	Method() string
}

type FileableConf interface {
	Configurable
	getFile() interface{}
	Url() string
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

func (cmc *CopyMessageConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(cmc.ChatID, 10))
	v.Add("from_chat_id", strconv.FormatInt(cmc.ChatID, 10))
	v.Add("message_id", strconv.FormatInt(cmc.MessageID, 10))
	// TODO: Make Optional fields too...
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
func (smc *SendMessageConfig) Values() (*url.Values, error) {
	result := &url.Values{}
	result.Add("chat_id", strconv.Itoa(smc.ChatID))
	result.Add("text", smc.Text)
	if smc.ParseMode != "" {
		result.Add("parse_mode", smc.ParseMode)
	}
	result.Add("disable_web_page_preview", strconv.FormatBool(smc.DisableWebPagePreview))
	return result, nil
}



// WebhookConfig uses for Using as arguemnt
// You may not fill all fields in struct
// https://core.telegram.org/bots/api#setwebhook
type SetWebhookConfig struct {
	URL            *url.URL // required
	Certificate    interface{}
	Offset         int
	MaxConnections int
	AllowedUpdates bool
	DropPendingUpdates bool
	IP string // if you need u can use it ;)
}

func (wc *SetWebhookConfig) Values() (*url.Values, error) {
	result := &url.Values{}
	// omg, why it s so bore ;(
	result.Add("url", wc.URL.String())
	// result.Add("certificate", wc.Certificate.URL())
	result.Add("ip_address", wc.IP)
	result.Add("max_connections", strconv.Itoa(wc.MaxConnections))
	result.Add("allowed_updates", strconv.FormatBool(wc.AllowedUpdates))
	result.Add("drop_pending_updates", strconv.FormatBool(wc.DropPendingUpdates))

	return result, nil
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

func (guc *GetUpdatesConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	if guc.Offset != 0 {
		v.Add("offset", strconv.Itoa(guc.Offset))
	}
	v.Add("limit", strconv.Itoa(guc.Limit))
	v.Add("timeout", strconv.Itoa(guc.Timeout))

	return v, nil
}

// Here methods name for various MethodConfigs
// 
func (cmc *CopyMessageConfig) Method() string { return "copyMessage" }
func (wc *SetWebhookConfig) Method() string { return "setWebhook" }
func (guc *GetUpdatesConfig) Method() string { return "getUpdates" }
func (smc *SendMessageConfig) Method() string { return "sendMessage" }
