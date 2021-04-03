package utils

import (
	"errors"
	"strings"

	"github.com/pikoUsername/tgp/objects"
)

// CheckToken Check out for a Space containing, and correct token
func CheckToken(token string) error {
	// Checks for space in token
	if strings.Contains(token, " ") {
		return errors.New("Token is Invalid! Token contains Space")
	}
	// Splits Token to 3 part as i know
	// Token contains 3 parts, first part is time creation
	// Second is i forget, 3 part is randomly generated part of token
	// Most inportant part of token
	token_parts := strings.Split(token, ":")
	if len(token_parts) > 3 {
		return errors.New("Token contains more than 3 parts")
	}
	// Checks for empty token
	if token_parts[0] == "" || token_parts[1] == "" || token_parts[2] == "" {
		return errors.New("Token is empty")
	}
	return nil
}

// Checks Statuscode and if Error then creates new Error with Error Description
func CheckResult(resp *objects.TelegramResponse) error {
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
