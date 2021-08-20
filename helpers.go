package tgp

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

var (
	DefaultContentTypes = NewContentTypes()
)

type ContentTypes struct {
	TEXT                              string
	AUDIO                             string
	DOCUMENT                          string
	ANIMATION                         string
	GAME                              string
	PHOTO                             string
	STICKER                           string
	VIDEO_NOTE                        string
	VOICE                             string
	CONTACT                           string
	LOCATION                          string
	VIDEO                             string
	VENUE                             string
	POLL                              string
	DICE                              string
	NEW_CHAT_MEMBERS                  string
	LEFT_CHAT_MEMBER                  string
	INVOICE                           string
	SUCCESSFUL_PAYMENT                string
	CONNECTED_WEBSITE                 string
	MESSAGE_AUTO_DELETE_TIMER_CHANGED string
	MIGRATE_TO_CHAT_ID                string
	MIGRATE_FROM_CHAT_ID              string
	PINNED_MESSAGE                    string
	NEW_CHAT_TITLE                    string
	NEW_CHAT_PHOTO                    string
	DELETE_CHAT_PHOTO                 string
	GROUP_CHAT_CREATED                string
	PASSPORT_DATA                     string
	PROXIMITY_ALERT_TRIGGERED         string
	VOICE_CHAT_SCHEDULED              string
	VOICE_CHAT_STARTED                string
	VOICE_CHAT_ENDED                  string
	VOICE_CHAT_PARTICIPANTS_INVITED   string

	UNKNOWN string
	ANY     string
}

func NewContentTypes() *ContentTypes {
	return &ContentTypes{
		TEXT:                              "TEXT",
		AUDIO:                             "AUDIO",
		DOCUMENT:                          "DOCUMENT",
		ANIMATION:                         "ANIMATION",
		GAME:                              "GAME",
		PHOTO:                             "PHOTO",
		STICKER:                           "STICKER",
		VIDEO:                             "VIDEO",
		VIDEO_NOTE:                        "VIDEO_NOTE",
		VOICE:                             "VOICE",
		CONTACT:                           "CONTACT",
		LOCATION:                          "LOCATION",
		VENUE:                             "VENUE",
		POLL:                              "POLL",
		DICE:                              "DICE",
		NEW_CHAT_MEMBERS:                  "NEW_CHAT_MEMBERS",
		LEFT_CHAT_MEMBER:                  "LEFT_CHAT_MEMBER",
		INVOICE:                           "INVOICE",
		SUCCESSFUL_PAYMENT:                "SUCESSFUL_PAYMENT",
		CONNECTED_WEBSITE:                 "CONNECTED_WEBSITE",
		MESSAGE_AUTO_DELETE_TIMER_CHANGED: "MESSAGE_AUTO_DELETE_TIMER_CHANGED",
		MIGRATE_TO_CHAT_ID:                "MIGRATE_TO_CHAT_ID",
		MIGRATE_FROM_CHAT_ID:              "MIGRATE_FROM_CHAT_ID",
		PINNED_MESSAGE:                    "PINNED_MESSAGE",
		NEW_CHAT_TITLE:                    "NEW_CHAT_TITLE",
		NEW_CHAT_PHOTO:                    "NEW_CHAT_PHOTO",
		DELETE_CHAT_PHOTO:                 "DELETE_CHAT_PHOTO",
		GROUP_CHAT_CREATED:                "GROUP_CHAT_CREATED",
		PASSPORT_DATA:                     "PASSPORT_DATA",
		PROXIMITY_ALERT_TRIGGERED:         "PROXIMITY_ALERT_TRIGGERED",
		VOICE_CHAT_SCHEDULED:              "VOICE_CHAT_SCHEDULED",
		VOICE_CHAT_STARTED:                "VOICE_CHAT_STARTED",
		VOICE_CHAT_ENDED:                  "VOICE_CHAT_ENDED",
		VOICE_CHAT_PARTICIPANTS_INVITED:   "VOICE_CHAT_PARTICIPANTS_INVITED",

		UNKNOWN: "UNKNOWN",
		ANY:     "ANY",
	}
}

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
	return (string)(marshal)
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

func getUidAndCidFromUpd(u *objects.Update) (cid_, uid_ int64) {
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

func fileToBytes(f interface{}, compress bool) ([]byte, error) {
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
			defer cw.Close()
			if err != nil {
				return []byte{}, err
			}
			if _, err := io.Copy(cw, fp); err != nil {
				return []byte{}, err
			}
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
			return "", errors.New("path is directory")
		}
		s = info.Name()
	default:
		return "", errors.New("reached, not reachable")
	}

	return s, nil
}

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
	var tgresp objects.TelegramResponse
	// Maybe use the Unmarshal ...
	err := json.NewDecoder(respBody).Decode(&tgresp)
	if err != nil {
		return &tgresp, err
	}
	return &tgresp, nil
}

// ReadFromInputFile call InputFile.Read method and for more shorter code this function was created
// using in Values() in SendDocumentConfig struct
func readFromInputFile(v *InputFile, compress bool) (p []byte, err error) {
	bs := make([]byte, v.Length)
	_, err = v.Read(bs)
	return bs, err
}

func ComposeFiles(files []*InputFile) { return }
