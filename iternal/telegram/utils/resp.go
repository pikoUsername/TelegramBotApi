package utils

import (
	"encoding/json"
	"io"

	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/ttypes"
)

// ResponseDecode ...
func ResponseDecode(respBody io.ReadCloser) (*ttypes.TelegramResponse, error) {
	var tgresp ttypes.TelegramResponse
	dec := json.NewDecoder(respBody)
	err := dec.Decode(&tgresp)
	if err != nil {
		return &tgresp, err
	}
	return &tgresp, nil
}
