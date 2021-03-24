package objects

// User represents Telegram User object
// docs: https://core.telegram.org/bots/api#user
type User struct {
	ID    int32 `json:"id"`
	IsBot bool  `json:"is_bot"`

	// Usernames
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`

	// You can use this, for i18n, or more for collect data from user ;(
	LanguageCode string `json:"language_code"`

	Location *ChatLocation `json:"location"`
}
