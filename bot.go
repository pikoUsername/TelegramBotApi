package tgp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pikoUsername/tgp/objects"
	"github.com/pikoUsername/tgp/utils"
	"github.com/technoweenie/multipartstreamer"
)

// HttpClient default interface for using by bot
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Bot can be created using Json config,
// Copy paste from go-telegram-bot-api ;D
type Bot struct {
	// Token uses for authonificate using URL
	// Url template {api_url}/bot{bot_token}/{method}?{args}
	Token string `json:"token"`

	// i will recomend to use HTML parse_mode
	// bc, HTML easy to use, and more conforatble
	ParseMode string `json:"parse_mode"`

	// Using prefix Bot, for avoid names conflict
	// and golang dont love name conflicts
	// by default this values is nil,
	// when you make get_me request, result
	// caches there, and you can take that
	// value in any moment.
	// Using Lazy method, instead of one moment
	Me *objects.User `json:"-"`

	// client if you need this, here
	// Client uses only for Post requests
	Client HttpClient `json:"-"`

	// default server must be here
	// if you wanna create own, just create
	// using this structure instead of NewBot function
	server *TelegramApiServer `json:"-"`

	// For DebugLog in console
	Debug bool `json:"debug"`
}

// NewBot get a new Bot
// This Fucntion checks a token
// for spaces and etc.
func NewBot(token string, checkToken bool, parseMode string) (*Bot, error) {
	if checkToken {
		// Check out for correct token
		err := utils.CheckToken(token)
		if err != nil {
			return nil, err
		}
	}
	return &Bot{
		Token:     token,
		ParseMode: parseMode,
		server:    DefaultTelegramServer,
		Client:    &http.Client{},
	}, nil
}

func (bot *Bot) Log(text string, v *url.Values, message interface{}) {
	if bot.Debug {
		log.Printf("%s req : %+v\n", text, v)
		log.Printf("%s resp: %+v\n", text, message)
	}
}

// MakeRequest to telegram servers
// and result parses to TelegramResponse
func (bot *Bot) MakeRequest(Method string, params *url.Values) (*objects.TelegramResponse, error) {
	// Bad Code, but working, huh

	// Creating URL
	// fix bug with sending request,
	// when url creates here or NewRequest not creates a correct url with url params
	tgurl := bot.server.ApiURL(bot.Token, Method)

	// Content Type is Application/json
	// Telegram uses application/json content type
	request, err := http.NewRequest("POST", tgurl, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	// Most important staff doing here
	// Sending Request to Telegram servers
	resp, err := bot.Client.Do(request)

	// check for error
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// make eatable
	tgresp, err := utils.ResponseDecode(resp.Body)
	if err != nil {
		return tgresp, err
	}
	return utils.CheckResult(tgresp)
}

// UploadFile same as MakeRequest, with one defference, file, and name variable, and nothing more
// copypaste of UploadFile go-telegram-bot
func (b *Bot) UploadFile(method string, f interface{}, fieldname string, values map[string]string) (*objects.TelegramResponse, error) {
	var name string
	ms := multipartstreamer.New()

	switch m := f.(type) {
	case string:
		name = m
		ms.WriteFile(fieldname, m)
	case InputFile:
		if m.URL != "" {
			values[fieldname] = m.URL

			ms.WriteFields(values)
		} else {
			data, err := ioutil.ReadAll(m.File)
			name = m.Name
			if name == "" {
				return &objects.TelegramResponse{}, errors.New("name field is nothing")
			}
			if err != nil {
				return &objects.TelegramResponse{}, err
			}

			buf := bytes.NewBuffer(data)

			ms.WriteReader(fieldname, m.Name, int64(len(data)), buf)
		}
	case FileableConf:
		name = m.Name()

		data, err := ioutil.ReadAll(m.GetFile().(io.Reader))
		if err != nil {
			return &objects.TelegramResponse{}, err
		}
		buf := bytes.NewBuffer(data)

		ms.WriteReader(fieldname, name, int64(len(data)), buf)
	case url.URL:
	case *url.URL:
		values[fieldname] = m.String()

		ms.WriteFields(values)
	default:
		return &objects.TelegramResponse{}, errors.New("not reached")
	}
	// creating File url
	tgurl := b.server.FileURL(b.Token, name)

	req, err := http.NewRequest(method, tgurl, nil)
	if err != nil {
		return &objects.TelegramResponse{}, err
	}
	ms.SetupRequest(req)

	// sending request
	resp, err := b.Client.Do(req)
	if err != nil {
		return &objects.TelegramResponse{}, err
	}
	// closing body
	b.Log("Response as bytes: ", nil, fmt.Sprintln(resp))
	defer resp.Body.Close()
	tgresp, err := utils.ResponseDecode(resp.Body)
	if err != nil {
		return tgresp, err
	}
	// returns response instant
	return utils.CheckResult(tgresp)
}

// GetMe reporesents telegram method
// https://core.telegram.org/bots/api#getme
func (bot *Bot) GetMe() (*objects.User, error) {
	if bot.Me != nil {
		return bot.Me, nil
	}
	resp, err := bot.MakeRequest("getMe", &url.Values{})
	if err != nil {
		return new(objects.User), err
	}
	var user objects.User
	err = json.Unmarshal(resp.Result, &user)
	if err != nil {
		return &user, err
	}

	bot.Me = &user // caching result
	return &user, nil
}

// Logout your bot from telegram
// https://core.telegram.org/bots/api#logout
func (bot *Bot) Logout() (*objects.TelegramResponse, error) {
	return bot.MakeRequest("logout", &url.Values{})
} // Indeed

// ===============================
// No returning value Telegram api methods
// ===============================

// DeleteChatPhoto represents deleteChatPhoto method
// https://core.telegram.org/bots/api#deletechatphoto
func (bot *Bot) DeleteChatPhoto(ChatId int64) (*objects.TelegramResponse, error) {
	v := &url.Values{}

	v.Add("chat_id", strconv.FormatInt(ChatId, 10))

	resp, err := bot.MakeRequest("deleteChatPhoto", v)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetchatTitle respresents setChatTitle method
// https://core.telegram.org/bots/api#setChatTitle
func (bot *Bot) SetChatTitle(ChatId int64, Title string) (*objects.TelegramResponse, error) {
	v := &url.Values{}

	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	v.Add("title", Title)

	resp, err := bot.MakeRequest("setChatTitle", v)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetChatDescription respresents setChatDescription method
// https://core.telegram.org/bots/api#setChatDescription
func (bot *Bot) SetChatDescription(ChatId int64, Description string) (*objects.TelegramResponse, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	v.Add("description", Description)
	resp, err := bot.MakeRequest("setChatDescription", v)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// PinChatMessage respresents pinChatMessage method
// https://core.telegram.org/bots/api#pinchatmessage
func (bot *Bot) PinChatMessage(
	ChatId int64,
	MessageId int64,
	DisableNotifiaction bool,
) (*objects.TelegramResponse, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	v.Add("message_id", strconv.FormatInt(MessageId, 10))
	v.Add("disable_notifications", strconv.FormatBool(DisableNotifiaction))
	resp, err := bot.MakeRequest("pinChatMessage", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UnpinAllChatMessage respresents unpinAllChatMessages method
// https://core.telegram.org/bots/api#unpinAllChatMessages
func (bot *Bot) UnpinAllChatMessages(ChatId int64) (*objects.TelegramResponse, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	resp, err := bot.MakeRequest("unpinAllChatMessages", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// =============================
// Message sending
// =============================

// Send uses as sender for almost all stuff
func (bot *Bot) SendMessageable(c Configurable) (*objects.Message, error) {
	v, err := c.Values()
	if err != nil {
		return &objects.Message{}, err
	}
	// Check out for parse_mode and set bot.ParseMode if config parse_mode is empty
	if v.Get("parse_mode") == "" {
		v.Set("parse_mode", bot.ParseMode)
	}
	resp, err := bot.MakeRequest(c.Method(), v)
	if err != nil {
		return &objects.Message{}, err
	}
	var msg objects.Message
	err = json.Unmarshal(resp.Result, &msg)
	bot.Log("SendMessageable function activated:", v, &msg)
	if err != nil {
		return &msg, err
	}
	return &msg, nil
}

// uploadAndSend will send a Message with a new file to Telegram.
func (bot *Bot) uploadAndSend(config FileableConf) (*objects.Message, error) {
	params, err := config.Params()
	if err != nil {
		return &objects.Message{}, err
	}

	file := config.GetFile()
	method := config.Method()
	resp, err := bot.UploadFile(method, file, config.Name(), params)
	if err != nil {
		return &objects.Message{}, err
	}

	var message *objects.Message
	json.Unmarshal(resp.Result, &message)

	bot.Log(method, nil, message)

	return message, nil
}

// Send ...
func (bot *Bot) Send(config Configurable) (*objects.Message, error) {
	switch c := config.(type) {
	case FileableConf:
		return bot.uploadAndSend(c)
	default:
		return bot.SendMessageable(c)
	}
}

// CopyMessage copies message
// https://core.telegram.org/bots/api#copymessage
func (bot *Bot) CopyMessage(config *CopyMessageConfig) (*objects.MessageID, error) {
	// Stub here, TODO: make for every config a values function/method
	v, err := config.Values()

	if err != nil {
		return &objects.MessageID{}, err
	}
	resp, err := bot.MakeRequest(config.Method(), v)
	if err != nil {
		return &objects.MessageID{}, err
	}
	var msg objects.MessageID

	err = json.Unmarshal(resp.Result, &msg)
	if err != nil {
		return &msg, err
	}
	return &msg, nil
}

// SendPhoto ...
func (bot *Bot) SendPhoto(config *SendPhotoConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendAudio ...
func (bot *Bot) SendAudio(config *SendAudioConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendDocument ...
func (bot *Bot) SendDocument(config *SendDocumentConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendVideo ...
func (bot *Bot) SendVideo(config *SendVideoConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendAnimation ...
func (bot *Bot) SendAnimation(config *SendAnimationConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendVoice ...
func (bot *Bot) SendVoice(config *SendVoiceConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendVideoName ...
func (bot *Bot) SendVideoName(config *SendVideoNoteConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendMediaGroup ...
func (bot *Bot) SendMediaGroup(config *SendMediaGroupConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendLocation ...
func (bot *Bot) SendLocation(config *SendLocationConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// editMessageLiveLocation ...
func (bot *Bot) EditMessageLiveLocation(config *EditMessageLLConf) (*objects.Message, error) {
	return bot.Send(config)
}

// SendMessage sends message using ChatID
// see: https://core.telegram.org/bots/api#sendmessage
func (bot *Bot) SendMessage(config *SendMessageConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendPoll Use this method to send a native poll
// https://core.telegram.org/bots/api#sendpoll
func (bot *Bot) SendPoll(config *SendPollConfig) (*objects.Message, error) {
	return bot.Send(config)
}

func (bot *Bot) SendDice(config *SendDiceConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// =========================
// Commands Methods
// =========================

// SetMyCommands Setup command to Telegram bot
// https://core.telegram.org/bots/api#setmycommands
func (bot *Bot) SetMyCommands(conf *SetMyCommandsConfig) (bool, error) {
	v, err := conf.Values()
	if err != nil {
		return false, err
	}
	resp, err := bot.MakeRequest(conf.Method(), v)
	if err != nil {
		return false, err
	}
	var ok bool
	err = json.Unmarshal(resp.Result, &ok)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// GetMyCommands get from bot commands command
// https://core.telegram.org/bots/api#getmycommands
func (bot *Bot) GetMyCommands(c *GetMyCommandsConfig) ([]objects.BotCommand, error) {
	v, _ := c.Values()
	resp, err := bot.MakeRequest(c.Method(), v)
	if err != nil {
		return nil, err
	}
	var cmds []objects.BotCommand
	err = json.Unmarshal(resp.Result, &cmds)
	if err != nil {
		return cmds, err
	}
	return cmds, nil
}

// ======================
// Getting Updates
// ======================

// DeleteWebhook if result is True, will be nil, if not so err
// https://core.telegram.org/bots/api#deletewebhook
func (bot *Bot) DeleteWebhook(c *DeleteWebhookConfig) error {
	v, err := c.Values()
	if err != nil {
		return err
	}
	_, err = bot.MakeRequest(c.Method(), v)
	if err != nil {
		return err
	}
	return nil
}

// GetUpdates uses for long polling
// https://core.telegram.org/bots/api#getupdates
func (bot *Bot) GetUpdates(c *GetUpdatesConfig) ([]*objects.Update, error) {
	v, err := c.Values()
	if err != nil {
		return nil, err
	}
	resp, err := bot.MakeRequest(c.Method(), v)
	if err != nil {
		return nil, &objects.TelegramApiError{
			Code:               resp.ErrorCode,
			Description:        resp.Description,
			ResponseParameters: objects.ResponseParameters{},
		}
	}
	var upd []*objects.Update
	err = json.Unmarshal(resp.Result, &upd)
	if err != nil {
		return upd, err
	}
	return upd, nil
}

// SetWebhook make subscribe to telegram events
// or sends to telegram a request for make
// Subscribe to specific IP, and when user
// sends a message to your bot, Telegram know
// Your bot IP and sends to your bot a Update
// https://core.telegram.org/bots/api#setwebhook
func (bot *Bot) SetWebhook(config *SetWebhookConfig) (*objects.TelegramResponse, error) {
	v, err := config.Values()
	if err != nil {
		return &objects.TelegramResponse{}, err
	}
	resp, err := bot.MakeRequest(config.Method(), v)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetWebhookInfo not require parametrs
// https://core.telegram.org/bots/api#getwebhookinfo
func (bot *Bot) GetWebhookInfo() (*objects.WebhookInfo, error) {
	resp, err := bot.MakeRequest("getWebhookInfo", &url.Values{})
	if err != nil {
		return &objects.WebhookInfo{}, err
	}
	var wi objects.WebhookInfo
	err = json.Unmarshal(resp.Result, &wi)
	if err != nil {
		return &wi, err
	}
	return &wi, nil
}

// =====================
// Chat methods
// =====================

// SendChatAction Resrpesents sendChatAction method
// https://core.telegram.org/bots/api#sendchataction
// Use this method when you need to tell the user that something is happening on the bot's side.
// The status is set for 5 seconds or less
// (when a message arrives from your bot, Telegram clients clear its typing status).
// Returns True on success.
func (bot *Bot) SendChatAction(c SendChatActionConf) (bool, error) {
	v, err := c.Values()
	if err != nil {
		return false, err
	}
	resp, err := bot.MakeRequest(c.Method(), v)
	if err != nil {
		return false, nil
	}
	var ok bool
	err = json.Unmarshal(resp.Result, &ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// DeleteChatStickerSet represents deleteChatStickerSet method
// https://core.telegram.org/bots/api#deletechatstickerset
func (bot *Bot) DeleteChatStickerSet(chat_id int64) (bool, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(chat_id, 10))
	resp, err := bot.MakeRequest("deleteChatStickerSet", v)
	if err != nil {
		return false, err
	}
	var ok bool
	err = json.Unmarshal(resp.Result, &ok)
	if err != nil {
		return ok, err
	}
	return ok, nil
}

// GetChat ...
func (bot *Bot) GetChat(chat_id int64) (*objects.Chat, error) {
	v := &url.Values{}
	v.Add("chat_id", strconv.FormatInt(chat_id, 10))

	resp, err := bot.MakeRequest("getChat", v)

	if err != nil {
		return &objects.Chat{}, err
	}

	var chat objects.Chat
	err = json.Unmarshal(resp.Result, &chat)

	if err != nil {
		return &chat, nil
	}

	return &chat, nil
}

// ================
// User methods
// ================

// GetUserProfilePhotos resresents getUserProfilePhotos method
// https://core.telegram.org/bots/api#getuserprofilephotos
func (bot *Bot) GetUserProfilePhotos(c GetUserProfilePhotosConf) (*objects.UserProfilePhotos, error) {
	v, _ := c.Values()
	resp, err := bot.MakeRequest(c.Method(), v)

	if err != nil {
		return &objects.UserProfilePhotos{}, nil
	}

	var photos objects.UserProfilePhotos
	err = json.Unmarshal(resp.Result, &photos)

	if err != nil {
		return &photos, err
	}

	return &photos, nil
}

// ====================
// other method
// ====================
