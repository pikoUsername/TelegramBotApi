package tgp

import (
	"errors"
	"strings"

	"github.com/pikoUsername/tgp/objects"
)

// CheckToken Check out for a Space containing, and correct token
func checkToken(token string) error {
	// Checks for space in token
	if strings.Contains(token, " ") {
		return errors.New("token is invalid! token contains space")
	}
	token_parts := strings.Split(token, ":")
	if len(token_parts) != 2 {
		return errors.New("token contains more than 2 parts")
	}
	// Checks for empty token
	if token_parts[0] == "" || token_parts[1] == "" {
		return errors.New("token is empty")
	}
	return nil
}

// Checks Statuscode and if Error then creates new Error with Error Description
func checkResult(resp *objects.TelegramResponse) (*objects.TelegramResponse, error) {
	// Check for Status, When StatusCode is 0 is default value
	// and Check is complete, and why so?
	// Telegram sends OK instead StatusCode 200
	if !resp.Ok {
		parameters := objects.ResponseParameters{}
		if resp.Parametrs != nil {
			parameters = *resp.Parametrs
		}
		return resp, &objects.TelegramApiError{
			Code:               resp.ErrorCode,
			Description:        resp.Description,
			ResponseParameters: parameters,
		}
	}

	return resp, nil
}
