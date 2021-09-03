package objects

import (
	"encoding/json"
	"unsafe"
)

type ReplyKBbuttonsType [][]KeyboardButton

// KeyboardButtonPoll represents keyboardButtonPollType object
// https://core.telegram.org/bots/api#keyboardbuttonpolltype
type KeyboardButtonPollType struct {
	Type string `json:"type"`
}

// KeyboardButton represents KeyboardButton object
// https://core.telegram.org/bots/api#keyboardbutton
type KeyboardButton struct {
	Text            string                  `json:"text"`
	RequestContact  bool                    `json:"request_contact"`
	RequestLocation bool                    `json:"request_location"`
	RequestPoll     *KeyboardButtonPollType `json:"request_poll"`
}

// ReplyKeyboardMarkup represents ReplyKeyboardMarkup object
// https://core.telegram.org/bots/api#replykeyboardmarkup
type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton
	ResizeKeyboard        bool   `json:"resize_keyboard"`
	OneTimeKeyboard       bool   `json:"one_time_keyboard"`
	InputFieldPlaceHolder string `json:"input_field_placeholder"`
	Selective             bool   `json:"selective"`
	RowWidth              uint
}

// Add ...
func (rkm *ReplyKeyboardMarkup) Add(btns ...KeyboardButton) *ReplyKeyboardMarkup {
	var row []KeyboardButton

	for index, btn := range btns {
		row = append(row, btn)
		if index%int(rkm.RowWidth) == 0 {
			rkm.Keyboard = append(rkm.Keyboard, row)
			row = []KeyboardButton{}
		}
	}
	if len(row) > 0 {
		rkm.Keyboard = append(rkm.Keyboard, row)
	}

	return rkm
}

func (rkm *ReplyKeyboardMarkup) String() string {
	v, _ := json.Marshal(rkm.Keyboard)
	return *(*string)(unsafe.Pointer(&v))
}
