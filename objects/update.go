package objects

// Update Represents telegram Update object
// https://core.telegram.org/bots/api#update
type Update struct {
	UpdateID           int `json:"update_id"`
	*Message           `json:"message"`
	EditedMessage      *Message            `json:"edited_message"`
	ChannelPost        *Message            `json:"channel_post"`
	EditedChannelPost  *Message            `json:"edited_channel_post"`
	InlineQuery        *InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *CallbackQuery      `json:"callback_query"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query"`
	*Poll              `json:"poll"`
	*PollAnswer        `json:"poll_answer"`
	MyChatMember       *ChatMemberUpdated `json:"my_chat_member"`
	ChatMember         *ChatMember        `json:"chat_member"`
	Date               int                `json:"date"`
	ForwardFrom        *User              `json:"forward_from"`
	ForwardDate        int                `json:"forward_date"`
	Dice               *Dice              `json:"dice"`

	// Using for data in handlers, NOT USING!
	DATA map[string]interface{}
}
