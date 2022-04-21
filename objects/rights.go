package objects

type ChatAdministratorRights struct {
	IsAnonymous         bool `json:"is_anonymous"`
	CanManageChat       bool `json:"can_manage_chat"`
	CanDeleteMessage    bool `json:"can_delete_message"`
	CanManageVoicechats bool `json:"can_manage_voice_chats"`
	CanRestrictMembers  bool `json:"can_restrict_members"`
	CanPromoteMembers   bool `json:"can_promote_members"`
	CanChangeInfo       bool `json:"can_change_info"`
	CanInviteUsers      bool `json:"can_invite_users"`
	CanPostMessages     bool `json:"can_post_message"`
	CanEditMessages     bool `json:"can_edit_messages"`
	CanPinMessages      bool `json:"can_pin_messages"`
}
