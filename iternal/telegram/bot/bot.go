package bot

import (
	"encoding/json"

	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/objects"
	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/utils"
)

// Bot ...
type Bot struct {
	Token string

	Me *objects.User `json:"-"`
}

// NewBot get a new Bot
// This Fucntion checks a token
// for spaces and etc.
func NewBot(token string, checkToken bool) (*Bot, error) {
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

// SendMessage
// see: https://core.telegram.org/bots/api#sendmessage
func (b *Bot) SendMessage(ChatID int, Text string) (*objects.Message, error) {
	// resp, err := MakeRequest("sendMessage", b.Token)
	return &objects.Message{}, nil
}
