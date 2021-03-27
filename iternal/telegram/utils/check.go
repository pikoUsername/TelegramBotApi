package utils

import (
	"strings"

	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/objects"
)

// CheckToken Check out for a Space containing
func CheckToken(token string) error {
	// Checks for space in token
	if strings.Contains(token, " ") {
		return InvalidToken
	}
	return nil
}

// Checks Statuscode and if Error then creates new Error with Error Description
func CheckResult(resp *objects.TelegramResponse) *error {
	// Check for Status, When StatusCode is 0 is default value
	// and Check is complete, and why so?
	// Telegram sends OK instead StatusCode 200
	if !resp.Ok {
		parameters := objects.ResponseParameters{}
		if resp.Parametrs != nil {
			parameters = *resp.Parametrs
		}
		return &objects.TelegramApiError{
			Code:               resp.ErrorCode,
			Description:        resp.Description,
			ResponseParameters: parameters,
		}
	}

	return nil
}
