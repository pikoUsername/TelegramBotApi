package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/pikoUsername/tgp/objects"
)

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

func ObjectToJson(obj interface{}) string {
	marshal, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func MarkupToString(obj interface{}) string {
	t, ok := obj.(*objects.InlineKeyboardMarkup)
	if ok {
		return ObjectToJson(t)
	}

	t2, ok := obj.(*objects.ReplyKeyboardMarkup)
	if ok {
		return ObjectToJson(t2)
	}

	return ""
}

func GetUidAndCidFromUpd(u *objects.Update) (cid_, uid_ int64) {
	var cid, uid int64

	if u.Message != nil {
		cid = u.Message.Chat.ID
		uid = u.Message.From.ID
	} else if u.EditedMessage != nil {
		cid = u.EditedMessage.Chat.ID
		uid = u.EditedMessage.From.ID
	} else if u.PinnedMessage != nil {
		cid = u.PinnedMessage.Chat.ID
		uid = u.PinnedMessage.From.ID
	} else if u.ReplyToMessage != nil {
		cid = u.ReplyToMessage.Chat.ID
		uid = u.ReplyToMessage.From.ID
	}

	return cid, uid
}

func FileToBytes(f interface{}, compress bool) ([]byte, error) {
	var bs []byte

	switch f := f.(type) {
	case string:
		fp, err := os.Open(f)
		if err != nil {
			return []byte{}, err
		}
		defer fp.Close()
		if compress {
			compressed := &bytes.Buffer{}
			cw, err := gzip.NewWriterLevel(compressed, gzip.BestCompression)
			if err != nil {
				return []byte{}, err
			}
			if _, err := io.Copy(cw, fp); err != nil {
				return []byte{}, err
			}
			cw.Close()
			return ioutil.ReadAll(compressed)
		}

		return ioutil.ReadAll(fp)
	case os.File:
	case *os.File:
		return ioutil.ReadAll(f)
	case []byte:
		return f, nil
	}

	return bs, nil
}

func RequestToUpdate(req *http.Request) (*objects.Update, error) {
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

func GuessFileName(f interface{}) (string, error) {
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
			return "", errors.New("file is directory")
		}
		s = info.Name()
	default:
		return "", errors.New("reached, not reachable")
	}

	return s, nil
}

func UrlValuesToMapString(v *url.Values, w map[string]string) {
	for key := range *v {
		value := v.Get(key)
		raw := w[key]
		if raw == "" {
			w[key] = value
		}
	}
}
