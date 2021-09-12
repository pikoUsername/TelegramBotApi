package tgp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unsafe"

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
	return BytesToString(marshal)
}

func FormatMarkup(obj interface{}) string {
	t, ok := obj.(*objects.InlineKeyboardMarkup)
	if ok {
		return t.String()
	}

	t2, ok := obj.(*objects.ReplyKeyboardMarkup)
	if ok {
		return t2.String()
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
	var tgresp objects.TelegramResponse
	// Maybe use the Unmarshal ...
	err := json.NewDecoder(respBody).Decode(&tgresp)
	if err != nil {
		return &tgresp, err
	}
	return &tgresp, nil
}

// CheckToken Check out for a Space containing, and token correct
func checkToken(token string) error {
	// Checks for space in token
	if strings.Contains(token, " ") {
		return tgpErr.New("token is invalid! token contains space")
	}
	token_parts := strings.Split(token, ":")
	if len(token_parts) != 2 {
		return tgpErr.New("token contains more than 2 parts")
	}
	// Checks for empty token
	if token_parts[0] == "" || token_parts[1] == "" {
		return tgpErr.New("token is empty")
	}
	return nil
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

// // Golang example that creates an http client that leverages a SOCKS5 proxy and a DialContext
// func NewClientFromEnv() (*http.Client, error) {
// 	proxyHost := os.Getenv("PROXY_HOST")

// 	baseDialer := &net.Dialer{
// 		Timeout:   30 * time.Second,
// 		KeepAlive: 30 * time.Second,
// 	}
// 	var dialContext DialContext

// 	if proxyHost != "" {
// 		dialSocksProxy, err := proxy.SOCKS5("tcp", proxyHost, nil, baseDialer)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "Error creating SOCKS5 proxy")
// 		}
// 		if contextDialer, ok := dialSocksProxy.(proxy.ContextDialer); ok {
// 			dialContext = contextDialer.DialContext
// 		} else {
// 			return nil, errors.New("Failed type assertion to DialContext")
// 		}
// 		logger.Debug("Using SOCKS5 proxy for http client",
// 			zap.String("host", proxyHost),
// 		)
// 	} else {
// 		dialContext = (baseDialer).DialContext
// 	}

// 	httpClient = newClient(dialContext)
// 	return httpClient, nil
// }

// func newClient(dialContext DialContext) *http.Client {
// 	return &http.Client{
// 		Transport: &http.Transport{
// 			Proxy:                 http.ProxyFromEnvironment,
// 			DialContext:           dialContext,
// 			MaxIdleConns:          10,
// 			IdleConnTimeout:       60 * time.Second,
// 			TLSHandshakeTimeout:   10 * time.Second,
// 			ExpectContinueTimeout: 1 * time.Second,
// 			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
// 		},
// 	}
// }
