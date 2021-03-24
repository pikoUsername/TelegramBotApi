package objects

// Represents Telegram ResponseParameters object
// https://core.telegram.org/bots/api#responseparameters
type ResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id"`
	RetryAfter      int   `json:"retry_after"`
}

// telegram api error
// can be raised when your request
// not correct
// see: https://github.com/TelegramBotAPI/errors
// official docs: https://core.telegram.org/api/errors
type TelegramApiError struct {
	Code        int
	Description string
	ResponseParameters
}

func (e *TelegramApiError) Error() string {
	return e.Description
}
