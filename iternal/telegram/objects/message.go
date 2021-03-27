package objects

// Message Telegram object
// https://core.telegram.org/bots/api#message
type Message struct {
	// MessageId ...
	MessageID int32 `json:"message_id"`

	// From User
	From *User `json:"from"`

	// Reply to Message, can be replied too
	ReplyToMessage *Message `json:"reply_to_message"`

	// ViaBot is Bot user, All Bots is Users
	ViaBot *User `json:"via_bot"`

	// EditDate int64 will work until end of Universe ;)
	EditDate int64 `json:"edit_date"`

	// MediaGroupId idk what is it, but docs say so, you need to make that
	MediaGroupId string `json:"media_group_id"`

	// AuthorSignature ...
	AuthorSignature string `json:"author_signature"`

	// Text the most important part of Message struct
	Text string `json:"text"`

	// PinnedMessage in 99% bot will be blocked by user if bot will ping user
	PinnedMessage *Message `json:"pinned_message"`

	// ReplyMarkup second most important thing
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup"`
}

// MessageEntity Uses in Message struct
// https://core.telegram.org/bots/api#messageentity
type MessageEntity struct {
	Type     string `json:"type"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
	URL      string `json:"url"`
	User     *User  `json:"user"`
	Language string `json:"language"`
}

// MessageID, idk why it s need
// https://core.telegram.org/bots/api#messageid
type MessageID struct {
	messageID int32 `json:"message_id"`
}
