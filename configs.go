package tgp

import (
	"io"
	"net/url"
	"strconv"

	"github.com/pikoUsername/tgp/objects"
	"github.com/pikoUsername/tgp/utils"
)

// This file stores ALL method configs

// functions which startswith New and method name
// uses for creating configs which stores ONLY required paramters

// Configurable is interface for using by method
type Configurable interface {
	Values() (*url.Values, error)
	Method() string
}

// FileableConf config using for Files storing
type FileableConf interface {
	Configurable
	Params() (map[string]string, error)
	Name() string
	Path() string
	GetFile() interface{}
}

// InputFile interaced by FileableConf
// Uses as Abstract level for real file
type InputFile struct {
	Name string
	URL  string
	File io.Reader
}

type BaseFile struct {
	FileID      string
	UseExisting bool
	MimeType    string
	FileSize    int
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
func (smc *SendMessageConfig) Values() (*url.Values, error) {
	result := &url.Values{}
	result.Add("chat_id", strconv.FormatInt(smc.ChatID, 10))

	result.Add("text", smc.Text)

	if smc.ParseMode != "" {
		result.Add("parse_mode", smc.ParseMode)
	}

	result.Add("reply_markup", utils.MarkupToString(smc.ReplyKeyboard))
	result.Add("disable_web_page_preview", strconv.FormatBool(smc.DisableWebPagePreview))
	// Must be work!
	result.Add("entities", utils.ObjectToJson(smc.Entities))

	return result, nil
}

func (smc *SendMessageConfig) Method() string {
	return "sendMessage"
}

func NewSendMessage(chat_id int64, text string) *SendMessageConfig {
	return &SendMessageConfig{
		ChatID: chat_id,
		Text:   text,
	}
}

// WebhookConfig uses for Using as arguemnt
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

func (wc *SetWebhookConfig) Values() (*url.Values, error) {
	result := &url.Values{}
	result.Add("url", wc.URL)

	if wc.Certificate == nil {
		cert, err := utils.FileToBytes(wc.Certificate, true)
		if err != nil {
			return &url.Values{}, err
		}
		result.Add("certificate", string(cert))
	}
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

func NewSetWebhook(url string) *SetWebhookConfig {
	return &SetWebhookConfig{
		URL: url,
	}
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
	ChatId          int64
	Audio           InputFile
	Caption         string
	ParseMode       string
	CaptionEntities []*objects.MessageEntity
	Duration        uint
	Performer       string
	Title           string
	Thumb           InputFile
}

func (sac *SendAudioConfig) Values() (*url.Values, error) {
	v := &url.Values{}

	v.Add("chat_id", strconv.FormatInt(sac.ChatId, 10))
	// Btw how???
	// v.Add("audio", )
	v.Add("caption", sac.Caption)
	if sac.ParseMode != "" {
		v.Add("parse_mode", sac.ParseMode)
	}
	v.Add("duration", strconv.FormatUint(uint64(sac.Duration), 10))
	v.Add("performer", sac.Performer)
	if sac.Title != "" {
		v.Add("title", sac.Title)
	}
	v.Add("caption_entities", utils.ObjectToJson(sac.CaptionEntities))
	return v, nil
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

	// Must be generic object, but for first time you can use InlineKeyboardMarkup
	// TODO
	ReplyMarkup *objects.InlineKeyboardMarkup
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

	v.Add("caption_entities", utils.ObjectToJson(svc.CaptionEntities))
	// TODO: reply Markup parsing function

	return v, nil
}

func (svc *SendVoiceConfig) Method() string {
	return "sendVoice"
}

type SendVideoNoteConfig struct {
}

func (svnc *SendVideoNoteConfig) Values() (*url.Values, error) {
	return &url.Values{}, nil
}

func (svnc *SendVideoNoteConfig) Method() string {
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

func (smgc *SendMediaGroupConfig) Values() (*url.Values, error) {
	v := &url.Values{}

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

func (smgc *SendMediaGroupConfig) Method() string {
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

func (slc *SendLocationConfig) Values() (*url.Values, error) {
	v := &url.Values{}

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

func (slc *SendLocationConfig) Method() string {
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
func (llc *EditMessageLLConf) Values() (*url.Values, error) {
	v := &url.Values{}
	return v, nil // stub
}

func (llc *EditMessageLLConf) Method() string {
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

// GetUpdate method fields
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

// Uses for default values for Sending updates
func NewGetUpdateConfig(Offset int) *GetUpdatesConfig {
	return &GetUpdatesConfig{
		Offset:  Offset,
		Limit:   20,
		Timeout: 5,
	}
}

// DeleteMyCommandsConfig ...
type DeleteMyCommandsConfig struct {
	Scope        objects.BotCommandScope // optional
	LanguageCode string                  // optional
}

func (dmcc *DeleteMyCommandsConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("scope", utils.ObjectToJson(dmcc.Scope))
	if dmcc.LanguageCode != "" {
		v.Add("language_code", dmcc.LanguageCode)
	}
	return v, nil
}

func (dmcc *DeleteMyCommandsConfig) Method() string {
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

func (smcc *SetMyCommandsConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("commands", utils.ObjectToJson(smcc.Commands))
	if smcc.LanguageCode != "" {
		v.Add("language_code", smcc.LanguageCode)
	}
	if smcc.Scope != nil {
		v.Add("scope", utils.ObjectToJson(smcc.Scope))
	}
	return v, nil
}

func (smcc *SetMyCommandsConfig) Method() string {
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

func (dwc *DeleteWebhookConfig) Values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("drop_pending_updates", strconv.FormatBool(dwc.DropPendingUpdates))
	return v, nil
}

func (dwc *DeleteWebhookConfig) Method() string {
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
		v.Add("explanation_entities", utils.ObjectToJson(spc.ExplnationEntites))
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

func (gc *GetChat) Values() (*url.Values, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(gc.ChatID, 10))
	return v, nil
}

func (gc *GetChat) Method() string {
	return "NONE"
}

// GetUserProfilePhotosConf represents getUserProfilePhotos method fields
// https://core.telegram.org/bots/api#getUserProfilePhotos
type GetUserProfilePhotosConf struct {
	UserId int64
	Offset int
	Limit  int
}

func (guppc *GetUserProfilePhotosConf) Values() (*url.Values, error) {
	v := &url.Values{}

	v.Add("user_id", strconv.FormatInt(guppc.UserId, 10))
	v.Add("offset", strconv.Itoa(guppc.Offset))
	v.Add("limit", strconv.Itoa(guppc.Limit))

	return v, nil
}

func (guppc *GetUserProfilePhotosConf) Method() string {
	return "getUserProfilePhotos"
}

type SendChatActionConf struct {
	ChatID int64
	Action string // see utils for actions type
}

func (scac *SendChatActionConf) Values() (*url.Values, error) {
	v := &url.Values{}

	v.Add("chat_id", strconv.FormatInt(scac.ChatID, 10))
	v.Add("action", scac.Action)

	return v, nil
}

func (scac *SendChatActionConf) Method() string {
	return "sendChatAction"
}
