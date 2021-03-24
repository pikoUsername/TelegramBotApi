package bot

import (
	"encoding/json"

	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/ttypes"
	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/utils"
)

// Bot ...
type Bot struct {
	Token string

	Me *ttypes.User `json:"-"`
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
func (bot *Bot) GetMe() (*ttypes.User, error) {
	resp, err := MakeRequest(utils.GETME, bot.Token, nil)
	if err != nil {
		return &ttypes.User{}, err
	}
	var user ttypes.User
	json.Unmarshal(resp.Result, &user)
	bot.Me = &user
	return &user, nil
}

// SendMessage
// see: https://core.telegram.org/bots/api#sendmessage
func (b *Bot) SendMessage(ChatID int, Text string) (*ttypes.Message, error) {
	// resp, err := MakeRequest("sendMessage", b.Token)
	return &ttypes.Message{}, nil
}
