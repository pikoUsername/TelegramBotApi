package tgp

import (
	"io"
	"net/url"
	"strconv"
	"time"

	"github.com/pikoUsername/tgp/objects"
)

// This file stores ALL method configs

// functions which startswith New and method name
// uses for creating configs which stores ONLY required paramters

// Configurable is interface for using by method
type Configurable interface {
	values() (url.Values, error)
	method() string
}

// FileableConf config using for Files storing
type FileableConf interface {
	Configurable
	params() (map[string]string, error)
	name() string
	path() string
	getFile() io.Reader
}

// InputFile interaced by FileableConf
// Uses as File holder
type InputFile struct {
	Name   string
	URL    string
	Length uint
	File   io.Reader
}

// Read ...
// No compress choice
func (f *InputFile) Read(p []byte) (n int, err error) {
	if f.File != nil && f.URL == "" {
		return f.File.Read(p)
	}
	bs, err := fileToBytes(f.Name, true)
	if err != nil {
		return 0, err
	}
	p = append(p[:0], bs...)
	return len(bs), nil
}

// BaseChat taken from go-telegram-bot-api
type BaseChat struct {
	ChatID              int64
	ChannelUsername     string
	ReplyToMessageID    int
	ReplyMarkup         interface{}
	DisableNotification bool
}

func (chat *BaseChat) params() (map[string]string, error) {
	params := make(map[string]string)

	return params, nil
}

// values returns url.Values representation of BaseChat
func (c *BaseChat) values() (url.Values, error) {
	v := url.Values{}
	if c.ChannelUsername != "" {
		v.Add("chat_id", c.ChannelUsername)
	} else {
		v.Add("chat_id", strconv.FormatInt(c.ChatID, 10))
	}

	if c.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(c.ReplyToMessageID))
	}

	if c.ReplyMarkup != nil {
		v.Add("reply_markup", MarkupToString(c.ReplyMarkup))
	}

	v.Add("disable_notification", strconv.FormatBool(c.DisableNotification))

	return v, nil
}

// BaseFile taken from go-telegram-bot-api
type BaseFile struct {
	BaseChat
	File        *InputFile
	FileID      string
	UseExisting bool
	MimeType    string
	FileSize    int
}

func (bf *BaseFile) params() (v map[string]string, err error) {
	v = make(map[string]string)

	if bf.FileID != "" {
		v["file_id"] = bf.FileID
	}
	v["use_existing"] = strconv.FormatBool(bf.UseExisting)
	if bf.MimeType != "" {
		v["mime_type"] = bf.MimeType
	}
	if bf.FileSize != 0 {
		v["file_size"] = strconv.Itoa(bf.FileSize)
	}

	cv, _ := bf.values()
	urlValuesToMapString(cv, v)

	return v, nil
}

func (bf *BaseFile) path() string {
	var path string
	if bf.File.Name != "" {
		path = bf.File.Name
	} else if bf.File.URL != "" {
		path = bf.File.URL
	}
	path = bf.FileID
	return path
}

func (bf *BaseFile) name() string {
	return bf.File.Name
}

func (bf *BaseFile) getFile() io.Reader {
	return bf.File
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

func (cmc *CopyMessageConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(cmc.ChatID, 10))
	v.Add("from_chat_id", strconv.FormatInt(cmc.ChatID, 10))
	v.Add("message_id", strconv.FormatInt(cmc.MessageID, 10))
	if cmc.Caption != "" {
		v.Add("caption", cmc.Caption)
	}
	if cmc.CaptionEntities != nil {
		v.Add("caption_entities", ObjectToJson(cmc.CaptionEntities))
	}
	v.Add("disable_notifications", strconv.FormatBool(cmc.DisableNotifications))
	if cmc.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(cmc.ReplyToMessageId, 10))
	}
	v.Add("allow_sending_with_reply", strconv.FormatBool(cmc.AllowSendingWithReply))
	if cmc.ReplyMarkup != nil {
		v.Add("reply_keyboards", MarkupToString(cmc.ReplyMarkup))
	}
	return v, nil
}

func (cmc *CopyMessageConfig) method() string {
	return "copyMessage"
}

// SendMessageConfig respresnests method,
// and fields of sendMessage method of telegram
// https://core.telegram.org/bots/api#sendmessage
type SendMessageConfig struct {
	// Required Field
	ChatID int64

	// It s too, Telegram excepts
	Text                  string // required
	ParseMode             string
	Entities              []*objects.MessageEntity
	DisableWebPagePreview bool

	DisableNotifiaction bool

	// ReplyKeyboard types:
	// InlineKeyboardMarkup
	// ReplyKeyboardMarkup
	ReplyKeyboard interface{}
}

// values ...
func (smc *SendMessageConfig) values() (url.Values, error) {
	result := url.Values{}
	result.Add("chat_id", strconv.FormatInt(smc.ChatID, 10))

	result.Add("text", smc.Text)

	if smc.ParseMode != "" {
		result.Add("parse_mode", smc.ParseMode)
	}

	result.Add("reply_markup", MarkupToString(smc.ReplyKeyboard))
	result.Add("disable_web_page_preview", strconv.FormatBool(smc.DisableWebPagePreview))
	// Must be work!
	result.Add("entities", ObjectToJson(smc.Entities))

	return result, nil
}

func (smc *SendMessageConfig) method() string {
	return "sendMessage"
}

func NewSendMessage(chat_id int64, text string) *SendMessageConfig {
	return &SendMessageConfig{
		ChatID: chat_id,
		Text:   text,
	}
}

// SetWebhookConfig uses for Using as arguemnt
// You may not fill all fields in struct
// https://core.telegram.org/bots/api#setwebhook
type SetWebhookConfig struct {
	URL                string // required
	Certificate        interface{}
	Offset             int
	MaxConnections     int
	AllowedUpdates     bool
	DropPendingUpdates bool
	IP                 string // if you need u can use it ;)
}

func (wc *SetWebhookConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("url", wc.URL)
	v.Add("ip_address", wc.IP) // required field
	if wc.MaxConnections != 0 {
		v.Add("max_connections", strconv.Itoa(wc.MaxConnections))
	}
	v.Add("allowed_updates", strconv.FormatBool(wc.AllowedUpdates))
	v.Add("drop_pending_updates", strconv.FormatBool(wc.DropPendingUpdates))

	return v, nil
}

func (wc *SetWebhookConfig) method() string {
	return "setWebhook"
}

func NewSetWebhook(url string) *SetWebhookConfig {
	return &SetWebhookConfig{
		URL: url,
	}
}

// SendPhotoConfig represnts telegram api method fields
// https://core.telegram.org/bots/api#sendphoto
type SendPhotoConfig struct {
	BaseFile
	Caption string
}

func (spc *SendPhotoConfig) values() (url.Values, error) {
	v, _ := spc.BaseChat.values()
	if spc.Caption != "" {
		v.Add("caption", spc.Caption)
	}
	return v, nil
}

func (spc *SendPhotoConfig) method() string {
	return "sendPhoto"
}

func NewSendPhoto(chat_id int64, photo *InputFile) *SendPhotoConfig {
	return &SendPhotoConfig{
		BaseFile: BaseFile{
			BaseChat: BaseChat{ChatID: chat_id},
			File:     photo,
		},
	}
}

// represents a sendAudio fields
type SendAudioConfig struct {
	BaseFile
	Caption         string
	ParseMode       string
	CaptionEntities []*objects.MessageEntity
	Duration        uint
	Performer       string
	Title           string
	// Thumb           *InputFile
}

func (sac *SendAudioConfig) values() (url.Values, error) {
	v, _ := sac.BaseFile.values()

	v.Add("chat_id", strconv.FormatInt(sac.ChatID, 10))

	if sac.Caption != "" {
		v.Add("caption", sac.Caption)
		if sac.ParseMode != "" {
			v.Add("parse_mode", sac.ParseMode)
		}
		if sac.CaptionEntities != nil {
			v.Add("caption_entities", ObjectToJson(sac.CaptionEntities))
		}
	}
	if sac.Duration != 0 {
		v.Add("duration", strconv.FormatUint(uint64(sac.Duration), 10))
	}
	if sac.Performer != "" {
		v.Add("performer", sac.Performer)
	}
	if sac.Title != "" {
		v.Add("title", sac.Title)
	}
	return v, nil
}

func (sac *SendAudioConfig) method() string {
	return "sendAudio"
}

func NewSendAudio(chatId int64, file *InputFile) *SendAudioConfig {
	return &SendAudioConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatId},
			File:        file,
			UseExisting: false,
		},
	}
}

// SendDocumentConfig represents sendDoucument method fields
type SendDocumentConfig struct {
	ChatID                      int64      // required
	Document                    *InputFile // required
	Thumb                       *InputFile
	Caption                     string
	ParseMode                   string
	CaptionEntities             []*objects.MessageEntity
	DisableContentTypeDetection bool
	DisableNotifiaction         bool
	ReplyToMessageID            int64
	AllowSendingWithoutReply    bool
	ReplyMarkup                 interface{}
}

func (sdc *SendDocumentConfig) values() (v url.Values, err error) {
	v = url.Values{}

	v.Add("chat_id", strconv.FormatInt(sdc.ChatID, 10))
	bs, err := readFromInputFile(sdc.Document, true)
	if err != nil {
		return v, err
	}
	v.Add("document", (string)(bs))
	if sdc.Caption != "" {
		v.Add("caption", sdc.Caption)
		if sdc.ParseMode != "" {
			v.Add("parse_mode", sdc.ParseMode)
		}
		if sdc.CaptionEntities != nil {
			v.Add("caption_entities", ObjectToJson(sdc.CaptionEntities))
		}
	}
	v.Add("disable_notification", strconv.FormatBool(sdc.DisableNotifiaction))
	if sdc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(sdc.ReplyToMessageID, 10))
	}
	v.Add("allow_sending_without_reply", strconv.FormatBool(sdc.AllowSendingWithoutReply))
	if sdc.ReplyMarkup != nil {
		v.Add("reply_markup", MarkupToString(sdc.ReplyMarkup))
	}
	return v, nil
}

func (sdc *SendDocumentConfig) name() string {
	return sdc.Document.Name
}

func (sdc *SendDocumentConfig) getFile() io.Reader {
	return sdc.Document
}

func (sdc *SendDocumentConfig) path() string {
	if sdc.Document != nil {
		return sdc.Document.Name
	}
	return sdc.Thumb.Name
}

func (sdc *SendDocumentConfig) params() (map[string]string, error) {
	params := make(map[string]string)

	v, _ := sdc.values()
	urlValuesToMapString(v, params)

	return params, nil
}

func (sdc *SendDocumentConfig) method() string {
	return "sendDocument"
}

func NewDocumentConfig(cid int64, r *InputFile) *SendDocumentConfig {
	return &SendDocumentConfig{
		ChatID:   cid,
		Document: r,
	}
}

// SendVideoConfig Represents sendVideo fields
// https://core.telegram.org/bots/api#sendvideo
type SendVideoConfig struct {
	BaseFile
	Duration uint32
	Width    uint16
	Height   uint16
	// Thumb    *InputFile
}

func (svc *SendVideoConfig) values() (url.Values, error) {
	v, _ := svc.BaseFile.values()
	if svc.Duration != 0 {
		v.Add("duration", strconv.FormatUint((uint64)(svc.Duration), 10))
	}

	if svc.Width != 0 {
		v.Add("width", strconv.FormatUint((uint64)(svc.Width), 10))
	}
	if svc.Height != 0 {
		v.Add("height", strconv.FormatUint((uint64)(svc.Height), 10))
	}

	return v, nil
}

func (svc *SendVideoConfig) method() string {
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

	// Thumb     *InputFile
	Caption   string
	ParseMode string
}

func (sac *SendAnimationConfig) values() (url.Values, error) {
	v := url.Values{}
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

func (sac *SendAnimationConfig) method() string {
	return "sendAnimation"
}

type SendVoiceConfig struct {
	BaseFile
	ChatId               int64
	Caption              string
	ParseMode            string
	CaptionEntities      []*objects.MessageEntity
	Duration             int
	DisableNotifications bool
	ReplyToMessageID     int64

	// Must be generic object, but for first time you can use InlineKeyboardMarkup
	// TODO
	ReplyMarkup *objects.InlineKeyboardMarkup
}

func (svc *SendVoiceConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(svc.ChatId, 10))
	v.Add("caption", svc.Caption)
	if svc.Caption != "" {
		v.Add("parse_mode", svc.Caption)
	}
	v.Add("disable_notifications", strconv.FormatBool(svc.DisableNotifications))
	if svc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(svc.ReplyToMessageID, 10))
	}

	v.Add("caption_entities", ObjectToJson(svc.CaptionEntities))
	// TODO: reply Markup parsing function

	return v, nil
}

func (svc *SendVoiceConfig) method() string {
	return "sendVoice"
}

type SendVideoNoteConfig struct {
	BaseFile
}

func (svnc *SendVideoNoteConfig) values() (url.Values, error) {
	return url.Values{}, nil
}

func (svnc *SendVideoNoteConfig) method() string {
	return "sendVideoName"
}

type SendMediaGroupConfig struct {
	// required fields
	ChatID int64
	Media  []interface{} // type: Union[[]InputMediaAudio, []InputMediaDocument, []InputMediaPhoto, []InputMediaVideo]

	// Optional fields
	DisableNotification      bool
	ReplyToMessageID         int64
	AllowSendingWithoutReply bool
}

func (smgc *SendMediaGroupConfig) values() (url.Values, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(smgc.ChatID, 10))
	// TOOD: media types
	// v.Add("media", smgc.Media)
	v.Add("disable_notification", strconv.FormatBool(smgc.DisableNotification))

	if smgc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(smgc.ReplyToMessageID, 10))
	}

	v.Add("allow_sending_without_reply", strconv.FormatBool(smgc.AllowSendingWithoutReply))

	return v, nil
}

func (smgc *SendMediaGroupConfig) method() string {
	return "sendMediaGroup"
}

func NewSendMediaGroupConfig(chat_id int64, media []interface{}) *SendMediaGroupConfig {
	return &SendMediaGroupConfig{
		ChatID: chat_id,
		Media:  media,
	}
}

type SendLocationConfig struct {
	ChatID                   int64   // req
	Latitude                 float32 // req
	Longitude                float32 // req
	HorizontalAccuracy       float32
	LivePeriod               uint
	Heading                  int
	ProximityAlertRadius     int
	DisableNotification      bool
	ReplyToMessageID         int
	AllowSendingWithoutReply bool
}

func (slc *SendLocationConfig) values() (url.Values, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(slc.ChatID, 10))

	// Same lines, broken DRY
	v.Add("latitude", strconv.FormatFloat(float64(slc.Latitude), 'E', -1, 64))
	v.Add("longitude", strconv.FormatFloat(float64(slc.Longitude), 'E', -1, 64))
	v.Add("horizontal_accuracy", strconv.FormatFloat(float64(slc.HorizontalAccuracy), 'E', -1, 64))

	if slc.LivePeriod != 0 {
		v.Add("live_period", strconv.FormatUint(uint64(slc.LivePeriod), 10))
	}

	if slc.Heading != 0 {
		v.Add("heading", strconv.FormatInt(int64(slc.Heading), 10))
	}

	v.Add("proximity_alert_radius", strconv.FormatInt(int64(slc.ProximityAlertRadius), 10))
	v.Add("disable_notification", strconv.FormatBool(slc.DisableNotification))

	if slc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(slc.ReplyToMessageID))
	}

	v.Add("allow_sending_without_reply", strconv.FormatBool(slc.AllowSendingWithoutReply))

	return v, nil
}

func NewSendLocationConf(chat_id int64, latitude float32, longitude float32) *SendLocationConfig {
	return &SendLocationConfig{
		ChatID:    chat_id,
		Latitude:  latitude,
		Longitude: longitude,
	}
}

func (slc *SendLocationConfig) method() string {
	return "sendLocation"
}

// LiveLocationConfig represents Telegram method fields of editmessageliveLocation
// https://core.telegram.org/bots/api#editmessagelivelocation
type EditMessageLLConf struct { // too long name anyway
	Longitude float32
	Latitude  float32
	ChatID    int64
	MessageID int64
}

// Values is stub!!
func (llc *EditMessageLLConf) values() (url.Values, error) {
	v := url.Values{}
	return v, nil // stub
}

func (llc *EditMessageLLConf) method() string {
	return "editMessageLiveLocation"
}

// all fields are required
func NewEditMessageLL(longitude float32, latit float32, chat_id int64, message_id int64) *EditMessageLLConf {
	return &EditMessageLLConf{
		Longitude: longitude,
		Latitude:  latit,
		ChatID:    chat_id,
		MessageID: message_id,
	}
}

type StopMessageLiveLocation struct {
	ChatID          int64
	MessageID       int64
	InlineMessageID int64
	ReplyMarkup     objects.InlineKeyboardMarkup
}

// GetUpdate method fields
type GetUpdatesConfig struct {
	Offset         int
	Limit          uint
	Timeout        uint
	AllowedUpdates []string
}

func (guc *GetUpdatesConfig) values() (url.Values, error) {
	v := url.Values{}
	if guc.Offset != 0 {
		v.Add("offset", strconv.Itoa(guc.Offset))
	}
	v.Add("limit", strconv.FormatUint(uint64(guc.Limit), 10))
	v.Add("timeout", strconv.FormatUint(uint64(guc.Timeout), 10))

	return v, nil
}

func (guc *GetUpdatesConfig) method() string {
	return "getUpdates"
}

// Uses for default values for Sending updates
func NewGetUpdateConfig(Offset int) *GetUpdatesConfig {
	return &GetUpdatesConfig{
		Offset:  Offset,
		Limit:   20,
		Timeout: 5,
	}
}

type GetMyCommandsConfig struct {
	Scope        objects.BotCommandScope // optional
	LanguageCode string                  // optional
}

func (gmcc *GetMyCommandsConfig) values() (url.Values, error) {
	v := url.Values{}
	if gmcc.Scope != nil {
		v.Add("scope", ObjectToJson(gmcc.Scope))
	}
	if gmcc.LanguageCode != "" {
		v.Add("language_code", gmcc.LanguageCode)
	}
	return v, nil
}

func (gmcc *GetMyCommandsConfig) method() string {
	return "getMyCommands"
}

// DeleteMyCommandsConfig ...
type DeleteMyCommandsConfig struct {
	Scope        objects.BotCommandScope // optional
	LanguageCode string                  // optional
}

func (dmcc *DeleteMyCommandsConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("scope", ObjectToJson(dmcc.Scope))
	if dmcc.LanguageCode != "" {
		v.Add("language_code", dmcc.LanguageCode)
	}
	return v, nil
}

func (dmcc *DeleteMyCommandsConfig) method() string {
	return "deleteMyCommands"
}

func NewDeleteMyCommandsConf() *DeleteMyCommandsConfig {
	return &DeleteMyCommandsConfig{}
}

// SetMyCommandsConfig ...
type SetMyCommandsConfig struct {
	Commands     []*objects.BotCommand
	Scope        objects.BotCommandScope
	LanguageCode string
}

func (smcc *SetMyCommandsConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("commands", ObjectToJson(smcc.Commands))
	if smcc.LanguageCode != "" {
		v.Add("language_code", smcc.LanguageCode)
	}
	if smcc.Scope != nil {
		v.Add("scope", ObjectToJson(smcc.Scope))
	}
	return v, nil
}

func (smcc *SetMyCommandsConfig) method() string {
	return "setMyCommands"
}

func NewSetMyCommands(commands ...*objects.BotCommand) *SetMyCommandsConfig {
	return &SetMyCommandsConfig{
		Commands: commands,
	}
}

// DeleteWebhookConfig ...
type DeleteWebhookConfig struct {
	DropPendingUpdates bool
}

func (dwc *DeleteWebhookConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("drop_pending_updates", strconv.FormatBool(dwc.DropPendingUpdates))
	return v, nil
}

func (dwc *DeleteWebhookConfig) method() string {
	return "deleteWebhook"
}

func NewDeleteWebHook(drop_pending_updates bool) *DeleteWebhookConfig {
	return &DeleteWebhookConfig{
		DropPendingUpdates: drop_pending_updates,
	}
}

// SendDiceConfig https://core.telegram.org/bots/api#senddice
type SendDiceConfig struct {
	ChatID                   int64
	Emoji                    string
	DisableNotifications     bool
	ReplyToMessageId         int64
	AllowSendingWithoutReply bool
	// ReplyMarkup will be type of objects.KeynoardMarkup not inline, and reply and etc.
	ReplyMarkup interface{}
}

func (sdc *SendDiceConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(sdc.ChatID, 10))
	if sdc.Emoji != "" {
		v.Add("emoji", sdc.Emoji)
	}
	v.Add("disable_notification", strconv.FormatBool(sdc.DisableNotifications))
	if sdc.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(sdc.ReplyToMessageId, 10))
	}
	v.Add("allow_sending_without_reply", strconv.FormatBool(sdc.AllowSendingWithoutReply))
	if sdc.ReplyMarkup != nil {
		v.Add("reply_markup", MarkupToString(sdc.ReplyMarkup))
	}
	return v, nil
}

func (sdc *SendDiceConfig) method() string {
	return "sendDice"
}

func NewSendDice(chatid int64, emoji string) *SendDiceConfig {
	return &SendDiceConfig{
		ChatID: chatid,
		Emoji:  emoji,
	}
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

func (spc *SendPollConfig) values() (url.Values, error) {
	v := url.Values{}
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
		v.Add("explanation_entities", ObjectToJson(spc.ExplnationEntites))
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

func (spc *SendPollConfig) method() string {
	return "sendPoll"
}

func NewSendPoll(chatid int64, question string, options []string) *SendPollConfig {
	return &SendPollConfig{
		ChatID:   chatid,
		Question: question,
		Options:  options,
	}
}

type GetChat struct {
	ChatID int64
}

func (gc *GetChat) values() (url.Values, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(gc.ChatID, 10))
	return v, nil
}

func (gc *GetChat) method() string {
	return "NONE"
}

// GetUserProfilePhotosConf represents getUserProfilePhotos method fields
// https://core.telegram.org/bots/api#getUserProfilePhotos
type GetUserProfilePhotosConf struct {
	UserId int64
	Offset int
	Limit  int
}

func (guppc *GetUserProfilePhotosConf) values() (url.Values, error) {
	v := url.Values{}

	v.Add("user_id", strconv.FormatInt(guppc.UserId, 10))
	v.Add("offset", strconv.Itoa(guppc.Offset))
	v.Add("limit", strconv.Itoa(guppc.Limit))

	return v, nil
}

func (guppc *GetUserProfilePhotosConf) method() string {
	return "getUserProfilePhotos"
}

type SendChatActionConf struct {
	ChatID int64
	Action string // see utils for actions type
}

func (scac *SendChatActionConf) values() (url.Values, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(scac.ChatID, 10))
	v.Add("action", scac.Action)

	return v, nil
}

func (scac *SendChatActionConf) method() string {
	return "sendChatAction"
}

type SendContactConfig struct {
	ChatID                   interface{} // req
	PhoneNumber              string      // req
	FirstName                string      // req
	LastName                 string
	Vcard                    string
	DisableNotifiaction      bool
	ReplyToMessageID         int64
	AllowSendingWithoutReply bool
	ReplyKeyboard            interface{}
}

func (scc *SendContactConfig) values() (url.Values, error) {
	v := url.Values{}
	switch t := scc.ChatID.(type) {
	case int64:
		v.Add("chat_id", strconv.FormatInt(t, 10))
	case string:
		v.Add("chat_id", t)
	}
	v.Add("phone_number", scc.PhoneNumber)
	v.Add("first_name", scc.FirstName)
	if scc.LastName != "" {
		v.Add("last_name", scc.LastName)
	}
	if scc.Vcard != "" {
		v.Add("vcard", scc.Vcard)
	}
	v.Add("disable_notification", strconv.FormatBool(!scc.DisableNotifiaction))
	if scc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(scc.ReplyToMessageID, 10))
	}
	if scc.ReplyKeyboard != nil {
		v.Add("reply_keyboard", MarkupToString(scc.ReplyKeyboard))
	}
	return v, nil
}

func (scc *SendContactConfig) method() string {
	return "sendContact"
}

// SendVenueConfig ...
type SendVenueConfig struct {
	ChatID                   interface{} // req
	Latitude                 float64     // req
	Longitude                float64     // req
	Title                    string      // req
	Address                  string      // req
	FoursQuareId             string
	FoursQuareType           string
	GooglePlaceId            string
	GooglePlaceType          string
	DisableNotification      bool
	ReplyToMessageId         int64
	AllowSendingWithoutReply bool
	ReplyMarkup              interface{}
}

func (svc *SendVenueConfig) values() (url.Values, error) {
	v := url.Values{}
	switch t := svc.ChatID.(type) {
	case int64:
		v.Add("chat_id", strconv.FormatInt(t, 10))
	case string:
		v.Add("chat_id", t)
	}
	v.Add("latitude", strconv.FormatFloat(svc.Latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(svc.Longitude, 'f', -1, 64))
	v.Add("title", svc.Title)
	v.Add("address", svc.Address)
	v.Add("allow_sending_without_reply", strconv.FormatBool(svc.AllowSendingWithoutReply))
	if svc.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(svc.ReplyToMessageId, 10))
	}
	v.Add("disable_notification", strconv.FormatBool(!svc.DisableNotification))
	if svc.GooglePlaceId != "" {
		v.Add("google_place_id", svc.GooglePlaceId)
	}
	if svc.GooglePlaceType != "" {
		v.Add("google_place_type", svc.GooglePlaceType)
	}
	if svc.FoursQuareId != "" {
		v.Add("four_square_id", svc.FoursQuareId)
	}
	if svc.FoursQuareType != "" {
		v.Add("four_square_type", svc.FoursQuareType)
	}
	if svc.ReplyMarkup != nil {
		v.Add("reply_markup", MarkupToString(svc.ReplyMarkup))
	}

	return v, nil
}

func (svc *SendVenueConfig) method() string {
	return "sendVenue"
}

// BanChatMemberConfig ...
type BanChatMemberConfig struct {
	ChatID         int64
	UserID         int64
	UntilDate      time.Duration
	RevokeMessages bool
}

func (bcm *BanChatMemberConfig) values() (url.Values, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(bcm.ChatID, 10))
	v.Add("user_id", strconv.FormatInt(bcm.UserID, 10))

	if bcm.UntilDate != 0 {
		v.Add("until_date", strconv.FormatInt((int64)(bcm.UntilDate), 10))
	}
	v.Add("revoke_messages", strconv.FormatBool(bcm.RevokeMessages))

	return v, nil
}

func (bcm *BanChatMemberConfig) method() string {
	return "banChatMember"
}

func NewBanChatMember(chat_id int64, user_id int64) *BanChatMemberConfig {
	return &BanChatMemberConfig{
		ChatID: chat_id,
		UserID: user_id,
	}
}
