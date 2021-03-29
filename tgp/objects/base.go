package objects

import (
	"encoding/json"
)

// MessageAble uses for Send function
// I love this, when you can instead just interface{}
// Use special interface, and your function can made
// This easily
type MessageAble interface {
	// Sorry but its impossible, make that without
	// circular imports, so i used this ;(
	// Huh, it just interface
	Send(interface{}) *Message
}

// maybe most important peace of code
type TelegramResponse struct {
	Ok          bool                `json:"ok"`
	Result      json.RawMessage     `json:"result"`
	ErrorCode   int                 `json:"error_code"`
	Description string              `json:"description"`
	Parametrs   *ResponseParameters `json:"parameters"`
}
