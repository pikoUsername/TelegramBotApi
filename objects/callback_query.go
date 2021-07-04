package objects

// CallbackQuery represents telegram object
// https://core.telegram.org/bots/api#callbackquery
type CallbackQuery struct {
	// Nothing to see, :P
	ID              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageID int64    `json:"inline_message_id"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`
	GameShortName   string   `json:"game_short_name"`
}
