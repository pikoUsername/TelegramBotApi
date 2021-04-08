package objects

// Chat type
type Chat struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

// https://core.telegram.org/bots/api#chatmember
type ChatMember struct {
	User        *User  `json:"user"`
	Status      string `json:"status"`
	CustomTitle string `json:"custom_title"`
	IsAnon      bool   `json:"is_anonymous"`
	IsMember    bool   `json:"is_member"`
	UntilDate   int64  `json:"until_date"`

	// Copy paste from perms, yeah broken DRY!
	CanBeEdited           bool `json:"can_be_edited"`
	CanManageChat         bool `json:"can_manage_chat"`
	CanPostMessages       bool `json:"can_post_message"`
	CanEditMessages       bool `json:"can_edit_messages"`
	CanDeleteMessage      bool `json:"can_delete_message"`
	CanManageVoicechats   bool `json:"can_manage_voice_chats"`
	CanRestrictMembers    bool `json:"can_restrict_members"`
	CanPromoteMembers     bool `json:"can_promote_members"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPinMessages        bool `json:"can_pin_messages"`
	CanSendMessage        bool `json:"can_send_message"`
	CanSendMediaMessages  bool `json:"can_send_media_messages"`
	CanSendPolls          bool `json:"can_send_polls"`
	CanSendOtherMessage   bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
}

type MyChatMember struct {
}

// ChatInviteLink represents ChatInvite object
// https://core.telegram.org/bots/api#chatinvitelink
type ChatInviteLink struct {
	InviteLink  string `json:"invite_link"`
	Creator     *User  `json:"creator"`
	IsPrimary   bool   `json:"is_primary"`
	IsRevoked   bool   `json:"is_revoked"`
	ExpireDate  int64  `json:"expire_date"`
	MemberLimit uint   `json:"member_limit"`
}

// ChatMemberUpdated object represents changes in the status of a chat member.
// https://core.telegram.org/bots/api#chatmemberupdated
type ChatMemberUpdated struct {
	Chat          *Chat           `json:"chat"`
	From          *User           `json:"user"`
	Date          uint64          `json:"date"`
	OldChatMember *ChatMember     `json:"old_chat_member"`
	NewChatMember *ChatMember     `json:"new_chat_member"`
	InviteLink    *ChatInviteLink `json:"invite_link"`
}
