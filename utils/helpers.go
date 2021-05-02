package utils

// Will be a other stuff, but useful stuff

// actions for telegram bot api
const (
	TYPING            = "TYPING"            // typing
	UPLOAD_PHOTO      = "UPLOAD_PHOTO"      // upload_photo
	RECORD_VIDEO      = "RECORD_VIDEO"      // record_video
	UPLOAD_VIDEO      = "UPLOAD_VIDEO"      // upload_video
	RECORD_AUDIO      = "RECORD_AUDIO"      // record_audio
	UPLOAD_AUDIO      = "UPLOAD_AUDIO"      // upload_audio
	RECORD_VOICE      = "RECORD_VOICE"      // record_voice
	UPLOAD_VOICE      = "UPLOAD_VOICE"      // upload_voice
	UPLOAD_DOCUMENT   = "UPLOAD_DOCUMENT"   // upload_docuemnt
	FIND_LOCATION     = "FIND_LOCATION"     // find_location
	RECORD_VIDEO_NOTE = "RECORD_VIDEO_NOTE" // record_video_note
	UPLOAD_VIDEO_NOTE = "UPLOAD_VIDEO_NOTE" // upload_video_note
)

// ParseModes
const (
	ModeMarkdown   = "Markdown"
	ModeMarkdownV2 = "MarkdownV2"
	ModeHTML       = "HTML"
)

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
