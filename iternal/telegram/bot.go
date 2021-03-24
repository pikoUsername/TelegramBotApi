package telegram

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/types"
	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/utils"
)

// Bot ...
type Bot struct {
	Token string

	Me *types.User `json:"-"`
}

// TelegramURL ...
var TelegramURL string = "https://api.telegram.org/"

// NewBot ...
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

// MakeRequest to telegram servers
// and result parses
func MakeRequest(Method string, Token string) (*types.TelegramResponse, error) {
	// Bad Code, but working, huh
	resp, err := http.Post(TelegramURL+"/"+Token+"/"+Method, "application/json", &strings.Reader{})
	if err != nil {
		return nil, errors.New("Error on sending request")
	}
	defer resp.Body.Close()

	var tgresp types.TelegramResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tgresp)
	if err != nil {
		return nil, err
	}

	if err := utils.CheckResult(Method, tgresp.ErrorCode, tgresp.Description); err != nil {
		return &tgresp, err
	}

	return &tgresp, nil
}

// GetMe ...
func (bot *Bot) GetMe() *types.User {
	return nil
}

// SendMessage ...
func (b *Bot) SendMessage(ChatID int, Text string) (*types.Message, error) {
	// resp, err := MakeRequest("SendMessage", b.Token)
	return &types.Message{}, nil
}
