package telegram

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Bot ...
type Bot struct {
	Token string
}

// TelegramURL ...
var TelegramURL string = "https://api.telegram.org/"

// MakeRequest to telegram servers
// and result parses
func MakeRequest(Method string, Token string) (*http.Response, error) {
	resp, err := http.Post(TelegramURL+"/"+Token+"/"+Method, "application/json", &strings.Reader{})
	if err != nil {
		return nil, errors.New("Error on sending request")
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(resp)
	if err != nil {
		return nil, errors.New("Error on json decode")
	}
	return resp, nil
}
