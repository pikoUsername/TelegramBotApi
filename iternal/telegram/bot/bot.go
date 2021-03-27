package bot

import (
	"encoding/json"
	"net/url"

	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/objects"
	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/utils"
)

// Bot ...
type Bot struct {
	Token string

	// Using prefix Bot, for avoid names conflict
	// and golang dont love name conflicts
	BotParseMode string
	Me           *objects.User `json:"-"`
}

// NewBot get a new Bot
// This Fucntion checks a token
// for spaces and etc.
func NewBot(token string, checkToken bool, parseMode string) (*Bot, error) {
	if checkToken {
		err := utils.CheckToken(token)
		if err != nil {
			return nil, err
		}
	}
	return &Bot{
		Token: token,
	}, nil
}

// GetMe reporesents telegram method
// https://core.telegram.org/bots/api#getme
func (bot *Bot) GetMe() (*objects.User, error) {
	resp, err := MakeRequest(GETME, bot.Token, nil)
	if err != nil {
		return &objects.User{}, err
	}
	var user objects.User
	json.Unmarshal(resp.Result, &user)
	bot.Me = &user
	return &user, nil
}

// Logout your bot from telegram
// https://core.telegram.org/bots/api#logout
func (bot *Bot) Logout() (bool, error) {
	_, err := MakeRequest(logout, bot.Token, url.Values{})
	if err != nil {
		return false, err
	}
	return true, nil
} // Indeed

// CopyMessage copies message
// https://core.telegram.org/bots/api#copymessage
func (bot *Bot) CopyMessage(config *CopyMessageConfig) (*objects.MessageID, error) {
	// Stub here, TODO: make for every config a values function/method
	resp, err := MakeRequest("copyMessage", bot.Token, url.Values{})
	if err != nil {
		return &objects.MessageID{}, err
	}
	var msg objects.MessageID

	json.Unmarshal(resp.Result, &msg)
	return &msg, nil
}

// GetUpdates uses for long polling
// https://core.telegram.org/bots/api#getupdates
func (bot *Bot) GetUpdates() (*objects.Update, error) {
	resp, err := MakeRequest(getUpdate, bot.Token, url.Values{})
	if err != nil {
		return &objects.Update{}, &objects.TelegramApiError{
			Code:               resp.ErrorCode,
			Description:        resp.Description,
			ResponseParameters: objects.ResponseParameters{},
		}
	}
	var upd objects.Update
	json.Unmarshal(resp.Result, &upd)
	return &upd, nil
}

// SendMessage sends message using ChatID
// see: https://core.telegram.org/bots/api#sendmessage
func (bot *Bot) SendMessage(ChatID int, Text string) (*objects.Message, error) {
	resp, err := MakeRequest("sendMessage", bot.Token, url.Values{})
	if err != nil {
		return &objects.Message{}, err
	}
	var msg objects.Message
	json.Unmarshal(resp.Result, &msg)
	return &msg, nil
}
