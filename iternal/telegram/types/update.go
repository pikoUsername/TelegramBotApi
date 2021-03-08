package types

// Update ...
type Update struct {
	UpdateID    int      `json:"update_id"`
	message     *Message `json:"message"`
	Date        int      `json:"date"`
	ForwardFrom *User    `json:"forward_from"`
	ForwardDate int      `json:"forward_date"`
	Dice        `json:"dice"`
}
