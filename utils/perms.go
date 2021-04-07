package utils

// Not Reporesenting any object
type Permissions struct {
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
