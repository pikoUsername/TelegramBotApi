package utils_test

import (
	"net/url"
	"testing"

	"github.com/pikoUsername/tgp/utils"
)

func TestURLToMapString(t *testing.T) {
	v := &url.Values{}
	base_value := "key"
	v.Add(base_value, base_value)
	params := make(map[string]string)
	utils.UrlValuesToMapString(v, params)
	b, ok := params[base_value]
	if !ok {
		t.Error("cannot get string by key")
		t.Fail()
	}

	if b != base_value {
		t.Error("values from params and real value is different")
		t.Fail()
	}
	t.Log(params)
}

const (
	TEXT                              = "TEXT"
	AUDIO                             = "AUDIO"
	DOCUMENT                          = "DOCUMENT"
	ANIMATION                         = "ANIMATION"
	GAME                              = "GAME"
	PHOTO                             = "PHOTO"
	STICKER                           = "STICKER"
	VIDEO                             = "VIDEO"
	VIDEO_NOTE                        = "VIDEO_NOTE"
	VOICE                             = "VOICE"
	CONTACT                           = "CONTACT"
	LOCATION                          = "LOCATION"
	VENUE                             = "VENUE"
	POLL                              = "POLL"
	DICE                              = "DICE"
	NEW_CHAT_MEMBERS                  = "NEW_CHAT_MEMBERS"
	LEFT_CHAT_MEMBER                  = "LEFT_CHAT_MEMBER"
	INVOICE                           = "INVOICE"
	SUCCESSFUL_PAYMENT                = "SUCESSFUL_PAYMENT"
	CONNECTED_WEBSITE                 = "CONNECTED_WEBSITE"
	MESSAGE_AUTO_DELETE_TIMER_CHANGED = "MESSAGE_AUTO_DELETE_TIMER_CHANGED"
	MIGRATE_TO_CHAT_ID                = "MIGRATE_TO_CHAT_ID"
	MIGRATE_FROM_CHAT_ID              = "MIGRATE_FROM_CHAT_ID"
	PINNED_MESSAGE                    = "PINNED_MESSAGE"
	NEW_CHAT_TITLE                    = "NEW_CHAT_TITLE"
	NEW_CHAT_PHOTO                    = "NEW_CHAT_PHOTO"
	DELETE_CHAT_PHOTO                 = "DELETE_CHAT_PHOTO"
	GROUP_CHAT_CREATED                = "GROUP_CHAT_CREATED"
	PASSPORT_DATA                     = "PASSPORT_DATA"
	PROXIMITY_ALERT_TRIGGERED         = "PROXIMITY_ALERT_TRIGGERED"
	VOICE_CHAT_SCHEDULED              = "VOICE_CHAT_SCHEDULED"
	VOICE_CHAT_STARTED                = "VOICE_CHAT_STARTED"
	VOICE_CHAT_ENDED                  = "VOICE_CHAT_ENDED"
	VOICE_CHAT_PARTICIPANTS_INVITED   = "VOICE_CHAT_PARTICIPANTS_INVITED"

	UNKNOWN = "UNKNOWN"
	ANY     = "ANY"
)
