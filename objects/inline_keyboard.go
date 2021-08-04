package objects

import "encoding/json"

// https://github.com/aiogram/aiogram/blob/dev-2.x/aiogram/types/inline_keyboard.py

// InlineKeyboardMarkup represents InlineKeyboardMarkup telegram object
// https://core.telegram.org/bots/api#inlinekeyboardmarkup
//
// Note: This will only work in Telegram versions released after
// 9 April, 2016. Older clients will display unsupported message.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
	RowWidth       uint
}

// Add ...
func (ikm *InlineKeyboardMarkup) Add(btns ...InlineKeyboardButton) *InlineKeyboardMarkup {
	var row []InlineKeyboardButton

	if ikm.RowWidth == 0 {
		ikm.RowWidth = 3
	}

	for index, btn := range btns {
		index += 1
		row = append(row, btn)
		if index%int(ikm.RowWidth) == 0 {
			ikm.InlineKeyboard = append(ikm.InlineKeyboard, row)
			row = []InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		ikm.InlineKeyboard = append(ikm.InlineKeyboard, row)
	}

	return ikm
}

func (self *InlineKeyboardMarkup) String() (string, error) {
	v, err := json.Marshal(self.InlineKeyboard)
	return (string)(v), err
}

// InlineKeyboardButton represnts InlineKeyboardButton telegram object
// https://core.telegram.org/bots/api#inlinekeyboardbutton
type InlineKeyboardButton struct {
	Text                         string        `json:"text"`
	URL                          string        `json:"url"`
	LoginURL                     *LoginURL     `json:"login_url"`
	CallbackData                 string        `json:"callback_data"`
	SwitchInlineQuery            string        `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string        `json:"switch_inline_query_current_chat"`
	CallbackGame                 *CallbackGame `json:"callback_game"`
	Pay                          bool          `json:"pay"`
}

func NewInlineKeyboardButton(text string, cd string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:         text,
		CallbackData: cd,
	}
}

func NewInlineKeyboardMarkup(row_width uint, btns ...InlineKeyboardButton) InlineKeyboardMarkup {
	ikm := InlineKeyboardMarkup{RowWidth: row_width}

	if len(btns) > 1 {
		ikm.Add(btns...)
	}

	return ikm
}
