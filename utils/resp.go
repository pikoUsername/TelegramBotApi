package utils

import (
	"encoding/json"
	"io"

	"github.com/pikoUsername/tgp/objects"
)

// ResponseDecode decodes to objects.TelegramResponse
// For next step parsing, in other function
// Result of Reponse saves in TelegramResponse.Result
func ResponseDecode(respBody io.ReadCloser) (*objects.TelegramResponse, error) {
	var tgresp objects.TelegramResponse
	// Maybe use the Unmarshal ...
	err := json.NewDecoder(respBody).Decode(&tgresp)
	if err != nil {
		return &tgresp, err
	}
	return &tgresp, nil
}
