package utils

import (
	"errors"
	"strings"
)

// CheckToken Check out for a Space containing
func CheckToken(token string) error {
	// Checks for space in token
	if strings.Contains(token, " ") {
		return errors.New("Token is Invalid! Token contains Space")
	}
	return nil
}

// Checks Statuscode and if Error then creates new Error with Error Description
func CheckResult(MethodName string, StatusCode int, Description string) error {
	// Check for Status, When StatusCode is 0 is default value
	// and Check is complete, and why so?
	// Telegram sends OK instead StatusCode 200
	if StatusCode != 0 || StatusCode == 200 {
		return errors.New(Description)
	}

	return nil
}
