package tg_types
// BotCommand respresents BotCommand object, ofc
// https://core.telegram.org/bots/api#botcommand
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}
