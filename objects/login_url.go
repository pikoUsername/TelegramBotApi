package objects

// LoginURL represents LoginURL telegram object
// https://core.telegram.org/bots/api#loginurl
type LoginURL struct {
	URL                string `json:"url"`
	ForwardText        string `json:"forward_text"`
	BotUsername        string `json:"bot_username"`
	RequestWriteAccess bool   `json:"request_write_access"`
}
