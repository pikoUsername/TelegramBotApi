package objects

import (
	"encoding/json"
)

// BaseFile is metadata for telegram file
// every telegram file object have this fields
type BaseFile struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FileName     string `json:"file_name"`
}

// TelegramResponse ...
type TelegramResponse struct {
	Ok          bool                `json:"ok"`
	Result      json.RawMessage     `json:"result"`
	ErrorCode   uint                `json:"error_code"`
	Description string              `json:"description"`
	Parametrs   *ResponseParameters `json:"parameters"`
}
