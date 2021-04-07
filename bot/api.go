package bot

import (
	"fmt"
)

type ITelegramServer interface {
	ApiURL(Token string, Method string) string
	FileURL(Token string, File string)
}

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
func (tas *TelegramApiServer) ApiURL(Token string, Method string) string {
	return fmt.Sprintf(tas.Base, Token, Method)
}

// FileUrl Creates at base of tas.File string
// a url for send a request
func (tas *TelegramApiServer) FileURL(Token string, File string) string {
	return fmt.Sprintf(tas.File, Token, File)
}

// Default telegram api server url
var (
	DefaultTelegramServer = NewTelegramApiServer("https://api.telegram.org")
)
