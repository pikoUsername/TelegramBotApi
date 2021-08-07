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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pikoUsername/tgp/objects"
	"github.com/technoweenie/multipartstreamer"
)

// StdLogger taken from logrus
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

// HttpClient default interface for using by bot
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Bot can be created using Json config,
// Copy pasted from go-telegram-bot-api ;D
//
// Client field have timeout, and default timeout is 5 second
// we can NOT change timeout if you using default client,
//
// Logger field is too default, and we can change logger any way we can
//
// ProxyURL is not used
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
	Me *objects.User `json:"me"`

	// client if you need this, here
	// Client uses only for Post requests
	Client HttpClient `json:"-"`

	// ProxyURL HTTP proxy URL
	ProxyURL *url.URL `json:"proxy_url"`

	// default server must be here
	// if you wanna create own, just create
	// using this structure instead of NewBot function
	server *TelegramApiServer

	// For DebugLog in console
	Debug bool `json:"debug"`

	// logger is one for dispatcher and Bot
	logger StdLogger `json:"-"`
}

// NewBot returns a New bot which need to contact with Telegram Bot API
// Bot structure should provide only Telegram bot API methods
func NewBot(token string, parseMode string, timeout time.Duration) (*Bot, error) {
	// Check out for correct token
	err := checkToken(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Token:     token,
		ParseMode: parseMode,
		server:    DefaultTelegramServer,
		logger:    log.New(os.Stderr, "", log.LstdFlags),
		// Client have default timeout 5 second
		Client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (bot *Bot) log(text string, v url.Values, message ...interface{}) {
	if bot.Debug {
		bot.logger.Printf("%s req : %+v\n", text, v)
		bot.logger.Printf("%s resp: %+v\n", text, message)
	}
}

// ===================
// sending requests
// ===================

// Request to telegram servers
// and result parses to TelegramResponse
func (bot *Bot) Request(Method string, params url.Values) (*objects.TelegramResponse, error) {
	// Bad Code, but working, huh

	// Creating URL
	// fix bug with sending request,
	// when url creates here or NewRequest not creates a correct url with url params
	tgurl := bot.server.ApiURL(bot.Token, Method+"?"+params.Encode())

	// Content Type is Application/json
	// Telegram uses application/json content type
	request, err := http.NewRequest("POST", tgurl, strings.NewReader(params.Encode()))
	if err != nil {
		return &objects.TelegramResponse{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	// Most important staff doing here
	// Sending Request to Telegram servers
	resp, err := bot.Client.Do(request)
	if err != nil {
		return &objects.TelegramResponse{}, err
	}
	defer resp.Body.Close()

	// make eatable
	tgresp, err := responseDecode(resp.Body)
	if err != nil {
		return tgresp, err
	}
	return checkResult(tgresp)
}

// DownloadFile uses for download file from any URL,
func (bot *Bot) DownloadFile(path string, w io.WriteSeeker, seek bool) error {
	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	resp, err := bot.Client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	if err != nil {
		return err
	}

	if seek {
		_, err := w.Seek((int64)(0), 0)
		if err != nil {
			return err
		}
	}
	return nil
}

// UploadFile same as MakeRequest, with one defference, file, and name variable, and nothing more
// copypaste of UploadFile go-telegram-bot
func (b *Bot) UploadFile(method string, f interface{}, fieldname string, v map[string]string) (*objects.TelegramResponse, error) {
	var name string
	ms := multipartstreamer.New()
	values := make(map[string]string)

	switch m := f.(type) {
	case string:
		name = m
		ms.WriteFile(fieldname, m)
	case *InputFile:
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
		name = m.name()

		data, err := ioutil.ReadAll(m.getFile())
		if err != nil {
			return &objects.TelegramResponse{}, err
		}
		buf := bytes.NewBuffer(data)

		ms.WriteReader(fieldname, name, (int64)(len(data)), buf)
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
	b.log("Response as bytes: ", nil, fmt.Sprintln(resp))
	defer resp.Body.Close()
	tgresp, err := responseDecode(resp.Body)
	if err != nil {
		return tgresp, err
	}
	// returns response instant
	return checkResult(tgresp)
}

// GetMe reporesents telegram method
// https://core.telegram.org/bots/api#getme
func (bot *Bot) GetMe() (*objects.User, error) {
	if bot.Me != nil {
		return bot.Me, nil
	}
	resp, err := bot.Request("getMe", url.Values{})
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
	return bot.Request("logout", url.Values{})
} // Indeed

// ===============================
// No returning value Telegram api methods
// ===============================

// DeleteChatPhoto represents deleteChatPhoto method
// https://core.telegram.org/bots/api#deletechatphoto
func (bot *Bot) DeleteChatPhoto(ChatId int64) (*objects.TelegramResponse, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(ChatId, 10))

	resp, err := bot.Request("deleteChatPhoto", v)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetchatTitle respresents setChatTitle method
// https://core.telegram.org/bots/api#setChatTitle
func (bot *Bot) SetChatTitle(ChatId int64, Title string) (*objects.TelegramResponse, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	v.Add("title", Title)

	resp, err := bot.Request("setChatTitle", v)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetChatDescription respresents setChatDescription method
// https://core.telegram.org/bots/api#setChatDescription
func (bot *Bot) SetChatDescription(ChatId int64, Description string) (*objects.TelegramResponse, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	v.Add("description", Description)
	resp, err := bot.Request("setChatDescription", v)
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
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	v.Add("message_id", strconv.FormatInt(MessageId, 10))
	v.Add("disable_notifications", strconv.FormatBool(DisableNotifiaction))
	resp, err := bot.Request("pinChatMessage", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UnpinAllChatMessage respresents unpinAllChatMessages method
// https://core.telegram.org/bots/api#unpinAllChatMessages
func (bot *Bot) UnpinAllChatMessages(ChatId int64) (*objects.TelegramResponse, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	resp, err := bot.Request("unpinAllChatMessages", v)
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
	v, err := c.values()
	if err != nil {
		return &objects.Message{}, err
	}
	// Check out for parse_mode and set bot.ParseMode if config parse_mode is empty
	if v.Get("parse_mode") == "" {
		v.Set("parse_mode", bot.ParseMode)
	}
	resp, err := bot.Request(c.method(), v)

	if err != nil {
		return &objects.Message{}, err
	}
	var msg objects.Message
	err = json.Unmarshal(resp.Result, &msg)
	bot.log("SendMessageable function activated:", v, &msg)
	if err != nil {
		return &msg, err
	}
	return &msg, nil
}

// uploadAndSend will send a Message with a new file to Telegram.
func (bot *Bot) UploadAndSend(config FileableConf) (*objects.Message, error) {
	params, err := config.params()
	if err != nil {
		return &objects.Message{}, err
	}

	file := config.getFile()
	method := config.method()
	resp, err := bot.UploadFile(method, file, config.name(), params)
	if err != nil {
		return &objects.Message{}, err
	}

	var message *objects.Message
	json.Unmarshal(resp.Result, &message)

	bot.log(method, nil, message)

	return message, nil
}

// Send ...
func (bot *Bot) Send(config interface{}) (*objects.Message, error) {
	switch c := config.(type) {
	case FileableConf:
		return bot.UploadAndSend(c)
	case Configurable:
		return bot.SendMessageable(c)
	}
	return &objects.Message{}, errors.New("config is not correct")
}

// CopyMessage copies message
// https://core.telegram.org/bots/api#copymessage
func (bot *Bot) CopyMessage(config *CopyMessageConfig) (*objects.MessageID, error) {
	v, err := config.values()

	if err != nil {
		return &objects.MessageID{}, err
	}
	resp, err := bot.Request(config.method(), v)
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

func (bot *Bot) SendContact(config *SendContactConfig) (*objects.Message, error) {
	return bot.Send(config)
}

func (bot *Bot) SendVenue(config *SendVenueConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// =========================
// Commands Methods
// =========================

// SetMyCommands Setup command to Telegram bot
// https://core.telegram.org/bots/api#setmycommands
func (bot *Bot) SetMyCommands(conf *SetMyCommandsConfig) (bool, error) {
	v, err := conf.values()
	if err != nil {
		return false, err
	}
	resp, err := bot.Request(conf.method(), v)
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
	v, _ := c.values()
	resp, err := bot.Request(c.method(), v)
	if err != nil {
		return []objects.BotCommand{}, err
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
func (bot *Bot) DeleteWebhook(c *DeleteWebhookConfig) (*objects.TelegramResponse, error) {
	v, err := c.values()
	if err != nil {
		return &objects.TelegramResponse{}, err
	}
	resp, err := bot.Request(c.method(), v)
	if err != nil {
		return &objects.TelegramResponse{}, err
	}
	return resp, nil
}

// GetUpdates uses for long polling
// https://core.telegram.org/bots/api#getupdates
func (bot *Bot) GetUpdates(c *GetUpdatesConfig) ([]*objects.Update, error) {
	v, err := c.values()
	if err != nil {
		return []*objects.Update{}, err
	}
	resp, err := bot.Request(c.method(), v)
	if err != nil {
		return []*objects.Update{}, &objects.TelegramApiError{
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
func (bot *Bot) SetWebhook(c *SetWebhookConfig) (*objects.TelegramResponse, error) {
	v, err := c.values()
	meth := c.method()

	// checkout for certificate, webhook may use without cert
	if c.Certificate == nil {
		return bot.Request(meth, v)
	}
	if err != nil {
		return &objects.TelegramResponse{}, err
	}
	params := make(map[string]string)
	urlValuesToMapString(v, params)
	// for debug
	bot.log("Params: ", nil, params)
	// uploads a certificate file, with other parametrs
	resp, err := bot.UploadFile(meth, c.Certificate, "certificate", params)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetWebhookInfo not require parametrs
// https://core.telegram.org/bots/api#getwebhookinfo
func (bot *Bot) GetWebhookInfo() (*objects.WebhookInfo, error) {
	resp, err := bot.Request("getWebhookInfo", url.Values{})
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
	v, err := c.values()
	if err != nil {
		return false, err
	}
	resp, err := bot.Request(c.method(), v)
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
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(chat_id, 10))
	resp, err := bot.Request("deleteChatStickerSet", v)
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
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(chat_id, 10))

	resp, err := bot.Request("getChat", v)

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
	v, _ := c.values()
	resp, err := bot.Request(c.method(), v)

	if err != nil {
		return &objects.UserProfilePhotos{}, err
	}

	var photos objects.UserProfilePhotos
	err = json.Unmarshal(resp.Result, &photos)

	if err != nil {
		return &photos, err
	}

	return &photos, nil
}

// ====================
// other methods
// ====================

func (bot *Bot) GetFile(file_id string) (*objects.File, error) {
	v := url.Values{}
	v.Add("file_id", file_id)
	resp, err := bot.Request("getFile", v)

	if err != nil {
		return &objects.File{}, err
	}

	var file *objects.File
	err = json.Unmarshal(resp.Result, &file)

	if err != nil {
		return file, err
	}

	return file, nil
}
