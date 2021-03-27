package objects

import (
	"encoding/json"
)

// maybe most important peace of code
type TelegramResponse struct {
	Ok          bool                `json:"ok"`
	Result      json.RawMessage     `json:"result"`
	ErrorCode   int                 `json:"error_code"`
	Description string              `json:"description"`
	Parametrs   *ResponseParameters `json:"parameters"`
}
