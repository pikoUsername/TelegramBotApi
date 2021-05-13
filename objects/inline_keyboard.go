package objects

// InlineKeyboardMarkup represents InlineKeyboardMarkup telegram object
// https://core.telegram.org/bots/api#inlinekeyboardmarkup
//
// Note: This will only work in Telegram versions released after
// 9 April, 2016. Older clients will display unsupported message.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton represnts InlineKeyboardButton telegram object
// https://core.telegram.org/bots/api#inlinekeyboardbutton
type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	URL                          string `json:"url"`
	*LoginURL                    `json:"login_url"`
	CallbackData                 string `json:"callback_data"`
	SwitchInlineQuery            string `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat"`
	*CallbackGame                `json:"callback_game"`
	Pay                          bool `json:"pay"`
}
