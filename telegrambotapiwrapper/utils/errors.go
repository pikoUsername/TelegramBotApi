package utils

import "errors"

var (
	InvalidToken error = errors.New("Token is Invalid! Token contains Space")
)
