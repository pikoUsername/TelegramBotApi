package bot

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pikoUsername/tgp/tgp/objects"
	"github.com/pikoUsername/tgp/tgp/utils"
)

// TelegramApiServer need in
// make easier use custom telegram api server
type TelegramApiServer struct {
	// Base telegram, sendMessage and etc.
	Base string

	// Url for file transfer, CDN and etc.
	File string
}

// NewTelegramApiServer ...
func NewTelegramApiServer(Base string) *TelegramApiServer {
	template := "/bot%s/%s"
	// /bot%s/%s means /bot<TOKEN>/<METHOD>
	return &TelegramApiServer{
		Base: fmt.Sprint(Base, template),
		File: fmt.Sprint(Base, "/file", template),
	}
}

// ApiUrl creates from base telegram url
func (tas *TelegramApiServer) ApiUrl(Token string, Method string) string {
	return fmt.Sprintf(tas.Base, Token, Method)
}

// FileUrl Creates at base of tas.File string
// a url for send a request
func (tas *TelegramApiServer) FileUrl(Token string, File string) string {
	return fmt.Sprintf(tas.File, Token, File)
}

// Default telegram api server url
var DefaultTelegramServer *TelegramApiServer = NewTelegramApiServer("https://api.telegram.org")

// MakeRequest to telegram servers
// and result parses to TelegramResponse
func MakeRequest(Method string, Token string, params *url.Values) (*objects.TelegramResponse, error) {
	// Bad Code, but working, huh
	// Content Type is Application/json
	// Telegram uses application/json content type
	cntype := "application/json"
	// Creating URL
	tgurl := DefaultTelegramServer.ApiUrl(Token, Method)

	resp, err := http.Post(tgurl, cntype, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// make eatable
	tgresp, err := utils.ResponseDecode(resp.Body)
	if err != nil {
		return tgresp, err
	}
	return tgresp, nil
}
