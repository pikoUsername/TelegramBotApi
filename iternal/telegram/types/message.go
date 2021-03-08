package types

// Message ...
type Message struct {
	MessageID int   `json:"message_id"`
	From      *User `json:"from"`
}
