package bot

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pikoUsername/tgp/configs"
	"github.com/pikoUsername/tgp/objects"
	"github.com/pikoUsername/tgp/utils"
)

// HttpClient ...
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Bot can be created using Json config,
// Copy paste from go-telegram-bot-api ;D
type Bot struct {
	Token     string `json:"token"`
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
	Debug bool
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

// GetMe reporesents telegram method
// https://core.telegram.org/bots/api#getme
func (bot *Bot) GetMe() (*objects.User, error) {
	if bot.Me != nil {
		return bot.Me, nil
	}
	resp, err := bot.MakeRequest("getMe", &url.Values{})
	if err != nil {
		return &objects.User{}, err
	}
	var user objects.User
	err = json.Unmarshal(resp.Result, &user)
	if err != nil {
		return &user, err
	}

	bot.Me = &user
	return &user, nil
}

// Logout your bot from telegram
// https://core.telegram.org/bots/api#logout
func (bot *Bot) Logout() (bool, error) {
	_, err := bot.MakeRequest("logout", &url.Values{})
	if err != nil {
		return false, err
	}
	return true, nil
} // Indeed

// =============================
// Message sending
// =============================

// Send uses as sender for almost all stuff
func (bot *Bot) SendMessageable(c configs.Configurable) (*objects.Message, error) {
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

// Send ...
func (bot *Bot) Send(config configs.Configurable) (*objects.Message, error) {
	switch config.(type) {
	case configs.FileableConf:
		return &objects.Message{}, nil
	default:
		return bot.SendMessageable(config)
	}
}

// CopyMessage copies message
// https://core.telegram.org/bots/api#copymessage
func (bot *Bot) CopyMessage(config *configs.CopyMessageConfig) (*objects.MessageID, error) {
	// Stub here, TODO: make for every config a values function/method
	v, err := config.Values()
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
func (bot *Bot) SendPhoto(config *configs.SendPhotoConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendAudio ...
func (bot *Bot) SendAudio(config *configs.SendAudioConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendDocument ...
func (bot *Bot) SendDocument(config *configs.SendDocumentConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendVideo ...
func (bot *Bot) SendVideo(config *configs.SendVideoConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendAnimation ...
func (bot *Bot) SendAnimation(config *configs.SendAnimationConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendVoice ...
func (bot *Bot) SendVoice(config *configs.SendVoiceConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendVideoName ...
func (bot *Bot) SendVideoName(config *configs.SendVideoNameConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendMediaGroup ...
func (bot *Bot) SendMediaGroup(config *configs.SendMediaGroupConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendLocation ...
func (bot *Bot) SendLocation(config *configs.SendLocationConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// editMessageLiveLocation ...
func (bot *Bot) EditMessageLiveLocation(config *configs.LiveLocationConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendMessage sends message using ChatID
// see: https://core.telegram.org/bots/api#sendmessage
func (bot *Bot) SendMessage(config *configs.SendMessageConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// =========================
// Commands Methods
// =========================

// SetMyCommands Setup command to Telegram bot
// https://core.telegram.org/bots/api#setmycommands
func (bot *Bot) SetMyCommands(conf *configs.SetMyCommandsConfig) (bool, error) {
	v, err := conf.Values() // Stub...
	if err != nil {
		return false, err
	}
	_, err = bot.MakeRequest(conf.Method(), v)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetMyCommands get from bot commands command
// https://core.telegram.org/bots/api#getmycommands
func (bot *Bot) GetMyCommands() ([]objects.BotCommand, error) {
	resp, err := bot.MakeRequest("getMyCommands", &url.Values{})
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
func (bot *Bot) DeleteWebhook(c *configs.DeleteWebhookConfig) error {
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
func (bot *Bot) GetUpdates(c *configs.GetUpdatesConfig) ([]objects.Update, error) {
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
	var upd []objects.Update
	log.Print(resp.Result)
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
func (bot *Bot) SetWebhook(config *configs.SetWebhookConfig) error {
	v, err := config.Values()
	if err != nil {
		return err
	}
	_, err = bot.MakeRequest(config.Method(), v)
	if err != nil {
		return err
	}
	return nil
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
