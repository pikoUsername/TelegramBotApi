package configs

import (
	"fmt"
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

// InputFile interaced by FileableConf
// Uses as Abstract level for real file
type InputFile struct {
	url  string
	File interface{}
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
	Text                  string // required
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
	result.Add("ip_address", wc.IP) // required field
	if wc.MaxConnections != 0 {
		result.Add("max_connections", strconv.Itoa(wc.MaxConnections))
	}
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

// SendVideoConfig Represents sendVideo fields
// https://core.telegram.org/bots/api#sendvideo
type SendVideoConfig struct {
	ChatId   int64
	Video    *InputFile
	Duration uint32
	Width    uint16
	Height   uint16
	Thumb    *InputFile
}

func (svc *SendVideoConfig) Values() (*url.Values, error) {
	return &url.Values{}, nil
}

func (svc *SendVideoConfig) Method() string {
	return "sendVideo"
}

// Represents Method SendAnimation Fields
// https://core.telegram.org/bots/api#sendanimation
type SendAnimationConfig struct {
	ChatId    int64      // ChatId might be a minus, or something like this
	Animation *InputFile // type: InputFile or string

	// Using unsigned, bc Duration width,
	// and height could be ONLY positive number
	Duration uint32
	Width    uint32 // Animation Width, what?
	Height   uint32

	Thumb     *InputFile
	Caption   string
	ParseMode string
}

func (sac *SendAnimationConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(sac.ChatId, 10))
	v.Add("duration", strconv.FormatUint(uint64(sac.Duration), 10))
	v.Add("width", strconv.FormatUint(uint64(sac.Width), 10))
	v.Add("height", strconv.FormatUint(uint64(sac.Height), 10))
	if sac.Caption != "" {
		v.Add("caption", sac.Caption)
	}
	if sac.ParseMode != "" {
		v.Add("parse_mode", sac.ParseMode)
	}

	return v, nil
}

func (sac *SendAnimationConfig) Method() string {
	return "sendAnimation"
}

type SendVoiceConfig struct {
	ChatId               int64
	Voice                interface{} // type: InputFile or String
	Caption              string
	ParseMode            string
	CaptionEntities      []*objects.MessageEntity
	Duration             int
	DisableNotifications bool
	ReplyToMessageID     int64
	ReplyMarkup          *objects.InlineKeyboard
}

func (svc *SendVoiceConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(svc.ChatId, 10))
	v.Add("caption", svc.Caption)
	if svc.Caption != "" {
		v.Add("parse_mode", svc.Caption)
	}
	v.Add("disable_notifications", strconv.FormatBool(svc.DisableNotifications))
	if svc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(svc.ReplyToMessageID, 10))
	}
	// TODO: reply Markup parsing function

	return v, nil
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
	Limit          uint
	Timeout        uint
	AllowedUpdates []string
}

func (guc *GetUpdatesConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	if guc.Offset != 0 {
		v.Add("offset", strconv.Itoa(guc.Offset))
	}
	v.Add("limit", strconv.FormatUint(uint64(guc.Limit), 10))
	v.Add("timeout", strconv.FormatUint(uint64(guc.Timeout), 10))

	return v, nil
}

func (guc *GetUpdatesConfig) Method() string {
	return "getUpdates"
}

type SetMyCommandsConfig struct {
	commands []*objects.BotCommand
}

func (smcc *SetMyCommandsConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("commands", "null") // Stub
	return v, nil
}

func (smcc *SetMyCommandsConfig) Method() string {
	return "setMyCommands"
}

type DeleteWebhookConfig struct {
	DropPendingUpdates bool
}

func (dwc *DeleteWebhookConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("drop_pending_updates", strconv.FormatBool(dwc.DropPendingUpdates))
	return v, nil
}

func (dwc *DeleteWebhookConfig) Method() string {
	return "deleteWebhook"
}

// SendDiceConfig https://core.telegram.org/bots/api#senddice
type SendDiceConfig struct {
	ChatID                   int64
	Emoji                    string
	DisableNotifications     bool
	ReplyToMessageId         int64
	AllowSendingWithoutReply bool
	// ReplyMarkup will be type of objects.KeynoardMarkup not inline, and reply and etc.
	// ReplyMarkup              *objects.KeyboardMarkup
}

func (sdc *SendDiceConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(sdc.ChatID, 10))
	if sdc.Emoji != "" {
		v.Add("emoji", sdc.Emoji)
	}
	v.Add("disable_notification", strconv.FormatBool(sdc.DisableNotifications))
	if sdc.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(sdc.ReplyToMessageId, 10))
	}
	v.Add("allow_sending_without_reply", strconv.FormatBool(sdc.AllowSendingWithoutReply))
	return v, nil
}

func (sdc *SendDiceConfig) Method() string {
	return "sendDice"
}

// SendPollConfig Use this method to send a native poll
// https://core.telegram.org/bots/api#sendpoll
type SendPollConfig struct {
	ChatID   int64
	Question string   // VarChar(300) limit 300 chars
	Options  []string // starts with 2->10 limit, 1-100 char limit

	// Vezet, Vezet
	IsAnonymous bool
	Type        string

	AllowsMultipleAnswers bool
	CorrectOptionId       int64
	Explanation           string
	ExpalnationParseMode  string
	ExplnationEntites     []*objects.MessageEntity

	// Using int time, here can be used time.Time
	OpenPeriod int64
	CloseDate  int64
	IsClosed   bool

	// Please, always turn off this
	DisableNotifications     bool
	ReplyToMessageID         int64
	AllowSendingWithoutReply bool
	// ReplyMarkup              *objects.KeyboardMarkup
}

func (spc *SendPollConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(spc.ChatID, 10))
	v.Add("question", spc.Question)
	// lucky, lucky
	v.Add("is_anonymous", strconv.FormatBool(spc.IsAnonymous))
	if spc.Type != "" {
		v.Add("type", spc.Type)
	}
	v.Add("allows_multiple_answers", strconv.FormatBool(spc.AllowsMultipleAnswers))
	v.Add("correct_option_id", strconv.FormatInt(spc.CorrectOptionId, 10))
	if spc.Explanation != "" {
		v.Add("explanation", spc.Explanation)
	}
	if spc.ExpalnationParseMode != "" {
		v.Add("explanation_parse_mode", spc.ExpalnationParseMode)
	}
	if spc.ExplnationEntites != nil {
		v.Add("explanation_entities", fmt.Sprintln(spc.ExplnationEntites))
	}
	v.Add("open_period", strconv.FormatInt(spc.OpenPeriod, 10))
	v.Add("close_date", strconv.FormatInt(spc.CloseDate, 10))
	v.Add("is_closed", strconv.FormatBool(spc.IsClosed))
	v.Add("disable_notifications", strconv.FormatBool(spc.DisableNotifications))
	if spc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(spc.ReplyToMessageID, 10))
	}
	return v, nil
}

func (spc *SendPollConfig) Method() string {
	return "sendPoll"
}
