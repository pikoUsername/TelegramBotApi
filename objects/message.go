package objects

// Message Telegram object(a huge object)
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
	SenderChat  *Chat                 `json:"sender_chat"`
	Chat        *Chat                 `json:"chat"`

	// Forwards
	ForwardFrom          *User  `json:"forward_from"`
	ForwardFromChat      *User  `json:"forward_from_chat"`
	ForwardFromMessageId int64  `json:"forward_from_message_id"`
	ForwardSignature     string `json:"forward_signature"`
	ForwardSenderName    string `json:"forward_sender_name"`
	ForwardDate          int64  `json:"forward_date"`

	// Start files fields
	Video *Video `json:"video"`

	// Hmmmm documentation...
	Document  *Document    `json:"document"`
	Animation *Animation   `json:"animation"`
	Photo     []*PhotoSize `json:"photo"`
	// Voice     *Voice       `json:"voice"`

	ConnectedWebsite string `json:"connected_website"`
	// Invoice          *Invoice `json:"invoice"`

	// Uses when user send message with photo
	Caption string `json:"caption"`

	// your location here
	*Location `json:"location"`

	// Chat stuff
	NewChatMembers []*ChatMember `json:"new_chat_members"`
	NewChatTitle   string        `json:"new_chat_title"`
	NewChatPhoto   []*PhotoSize  `json:"new_chat_photo"`
	LeftChatMember *User         `json:"left_chat_member"`
}

func (m *Message) GetContentType() string {
	// ContentTypes from utils/ can be used bc cycle import
	if m.Text != "" {
		return "TEXT"
	} else if m.Animation != nil {
		return "ANIMATION"
	} else {
		return "UNKNOWN"
	}
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
	MessageID int32 `json:"message_id"`
}
