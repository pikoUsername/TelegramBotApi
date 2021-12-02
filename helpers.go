package tgp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"unsafe"

	"github.com/pikoUsername/tgp/objects"
)

// Will be a other stuff, but useful stuff

// actions for sendAction method
const (
	TYPING            = "TYPING"
	UPLOAD_PHOTO      = "UPLOAD_PHOTO"
	RECORD_VIDEO      = "RECORD_VIDEO"
	UPLOAD_VIDEO      = "UPLOAD_VIDEO"
	RECORD_AUDIO      = "RECORD_AUDIO"
	UPLOAD_AUDIO      = "UPLOAD_AUDIO"
	RECORD_VOICE      = "RECORD_VOICE"
	UPLOAD_VOICE      = "UPLOAD_VOICE"
	UPLOAD_DOCUMENT   = "UPLOAD_DOCUMENT"
	FIND_LOCATION     = "FIND_LOCATION"
	RECORD_VIDEO_NOTE = "RECORD_VIDEO_NOTE"
	UPLOAD_VIDEO_NOTE = "UPLOAD_VIDEO_NOTE"
	CHOOSE_STICKER    = "CHOOSE_STICKER"
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

func ObjectToJson(obj interface{}) string {
	marshal, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return BytesToString(marshal)
}

func FormatMarkup(obj interface{}) string {
	switch t := obj.(type) {
	case *objects.ReplyKeyboardMarkup:
	case *objects.InlineKeyboardMarkup:
		return t.String()
	default:
	}

	return ""
}

func extractIds(u *objects.Update) (cid_, uid_ int64) {
	var cid, uid int64

	if u.Message != nil {
		cid = u.Message.Chat.ID
		uid = u.Message.From.ID
	} else if u.EditedMessage != nil {
		cid = u.EditedMessage.Chat.ID
		uid = u.EditedMessage.From.ID
	} else if u.Message.PinnedMessage != nil {
		cid = u.Message.PinnedMessage.Chat.ID
		uid = u.Message.PinnedMessage.From.ID
	} else if u.Message.ReplyToMessage != nil {
		cid = u.Message.ReplyToMessage.Chat.ID
		uid = u.Message.ReplyToMessage.From.ID
	}

	return cid, uid
}

func requestToUpdate(req *http.Request) (*objects.Update, error) {
	if req.Method != http.MethodPost {
		return &objects.Update{}, errors.New("wrong HTTP method required POST")
	}

	var update *objects.Update
	err := json.NewDecoder(req.Body).Decode(&update)
	if err != nil {
		return &objects.Update{}, err
	}

	return update, nil
}

func guessFileName(f interface{}) (string, error) {
	var s string

	switch f := f.(type) {
	case os.FileInfo:
		return f.Name(), nil

	case string:
		return f, nil

	case os.File:
		info, err := f.Stat()
		if err != nil {
			return "", err
		}
		if info.IsDir() {
			return "", tgpErr.New("path is directory")
		}
		s = info.Name()

	default:
	}

	return s, tgpErr.New("incorrect object type A , type must be in os.File, tgp.InputFile, string, os.FileInfo")
}

// Just write to map from url.Values
func urlValuesToMapString(v url.Values, w map[string]string) {
	for key, value := range v {
		if len(v[key]) > 0 {
			w[key] = value[0]
		}
	}
}

// ResponseDecode decodes to objects.TelegramResponse
// For next step parsing, in other function
// Result of Reponse saves in TelegramResponse.Result
func responseDecode(respBody io.ReadCloser) (*objects.TelegramResponse, error) {
	var tgresp *objects.TelegramResponse
	// Maybe use the Unmarshal ...
	err := json.NewDecoder(respBody).Decode(&tgresp)
	if err != nil {
		return tgresp, err
	}
	return tgresp, nil
}

// Checks Statuscode and if Error then creates new Error with Error Description
func checkResult(resp *objects.TelegramResponse) (*objects.TelegramResponse, error) {
	// Check for Status, When StatusCode is 0 is default value
	// and Check is complete, and why so?
	// Telegram sends OK instead StatusCode 200
	if !resp.Ok {
		parameters := objects.ResponseParameters{}
		if resp.Parametrs != nil {
			parameters = *resp.Parametrs
		}
		return resp, &objects.TelegramApiError{
			Code:               resp.ErrorCode,
			Description:        resp.Description,
			ResponseParameters: parameters,
		}
	}

	return resp, nil
}

// From gin-gonic/internal/bytesconv/bytesconv.go

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(x string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{x, len(x)},
	))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func WriteRequestError(wr http.ResponseWriter, err error) {
	errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
	wr.WriteHeader(http.StatusBadRequest)
	wr.Header().Set("Content-Type", "application/json")
	wr.Write(errMsg)
}
