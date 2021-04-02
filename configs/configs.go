package configs

import (
	"net/url"
	"strconv"

	"github.com/pikoUsername/tgp/objects"
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

func (cmc *CopyMessageConfig) Method() string {
	return "copyMessage"
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

func (smc *SendMessageConfig) Method() string {
	return "sendMessage"
}

// WebhookConfig uses for Using as arguemnt
// You may not fill all fields in struct
// https://core.telegram.org/bots/api#setwebhook
type SetWebhookConfig struct {
	URL                *url.URL // required
	Certificate        interface{}
	Offset             int
	MaxConnections     int
	AllowedUpdates     bool
	DropPendingUpdates bool
	IP                 string // if you need u can use it ;)
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

func (wc *SetWebhookConfig) Method() string {
	return "setWebhook"
}

// Here starts a stubs
// SendPhotoConfig ...
type SendPhotoConfig struct {
}

func (spc *SendPhotoConfig) Values() (*url.Values, error) {
	return &url.Values{}, nil
}

func (spc *SendPhotoConfig) Method() string {
	return "sendPhoto"
}

// represents a sendAudio fields
type SendAudioConfig struct {
}

func (sac *SendAudioConfig) Values() (*url.Values, error) {
	return &url.Values{}, nil
}

func (sac *SendAudioConfig) Method() string {
	return "sendAudio"
}

// SendDocumentConfig represents sendDoucument method fields
type SendDocumentConfig struct {
}

func (sdc *SendDocumentConfig) Values() (*url.Values, error) {
	return &url.Values{}, nil
}

func (sdc *SendDocumentConfig) Method() string {
	return "sendDocument"
}

type SendVideoConfig struct {
}

func (svc *SendVideoConfig) Values() (*url.Values, error) {
	return &url.Values{}, nil
}

func (svc *SendVideoConfig) Method() string {
	return "sendVideo"
}

type SendAnimationConfig struct {
}

func (sac *SendAnimationConfig) Values() (*url.Values, error) {
	return &url.Values{}, nil
}

func (sac *SendAnimationConfig) Method() string {
	return "sendAnimation"
}

type SendVoiceConfig struct {
}

func (svc *SendVoiceConfig) Values() (*url.Values, error) {
	return &url.Values{}, nil
}

func (svc *SendVoiceConfig) Method() string {
	return "sendVoice"
}

type SendVideoNameConfig struct {
}

func (svnc *SendVideoNameConfig) Values() (*url.Values, error) {
	return &url.Values{}, nil
}

func (svnc *SendVideoNameConfig) Method() string {
	return "sendVideoName"
}

type SendMediaGroupConfig struct {
}

func (smgc *SendMediaGroupConfig) Values() (*url.Values, error) {
	return &url.Values{}, nil
}

func (smgc *SendMediaGroupConfig) Method() string {
	return "sendMediaGroup"
}

type SendLocationConfig struct {
}

func (slc *SendLocationConfig) Values() (*url.Values, error) {
	return &url.Values{}, nil
}

func (slc *SendLocationConfig) Method() string {
	return "sendLocation"
}

// LiveLocationConfig represents Telegram method fields of liveLocation
// https://core.telegram.org/bots/api#editmessagelivelocation
type LiveLocationConfig struct {
	Longitude float32
	Latitude  float32
	ChatID    int64
	MessageID int64
}

// Values is stub!!
func (llc *LiveLocationConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	return v, nil // stub
}

func (llc *LiveLocationConfig) Method() string {
	return "editMessageLiveLocation"
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

func (guc *GetUpdatesConfig) Method() string {
	return "getUpdates"
}

// Here methods name for various Metho Configs
