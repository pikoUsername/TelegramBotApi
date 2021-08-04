package objects

import (
	"errors"
	"fmt"
)

// Represents Telegram ResponseParameters object
// https://core.telegram.org/bots/api#responseparameters
type ResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id"`
	RetryAfter      int   `json:"retry_after"`
}

// telegram api error
// can be raised when your request
// not correct
// see: https://github.com/TelegramBotAPI/errors
// official docs: https://core.telegram.org/api/errors
type TelegramApiError struct {
	Code        uint
	Description string
	ResponseParameters
}

func (e *TelegramApiError) Error() string {
	return fmt.Sprintf("telegram: %s", e.Description)
}

// ErrorPrefix get to user/client
// a error with prefix and splited up with separator
// used in errors variable, lol
type ErrorPrefix struct {
	prefix    string
	separator string
}

func (eg *ErrorPrefix) New(text string) error {
	return errors.New(eg.prefix + eg.separator + text)
}

func NewErrorPrefix(prefix, sep string) *ErrorPrefix {
	return &ErrorPrefix{prefix: prefix, separator: sep}
}

var (
	Errors = NewErrorPrefix("tgp", ": ")
)
