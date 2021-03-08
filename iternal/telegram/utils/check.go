package utils

import (
	"errors"
	"strings"
)

// CheckToken Check out for a Space containing
func CheckToken(token string) error {
	if strings.Contains(token, " ") {
		return errors.New("Token is Invalid! Token contains Space")
	}
	return nil
}
