package utils

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
