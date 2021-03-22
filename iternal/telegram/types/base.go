package types

import (
	"encoding/json"
)

type TelegramResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   string          `json:"error_code"`
	Description string          `json:"description"`
}
