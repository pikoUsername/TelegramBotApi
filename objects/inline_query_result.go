package objects

// dont scroll down

// unfotunatly package version is lower than 1.18. So there is no generics ;((((
type InlineQueryResult interface {
	getQueryResult() // stub function
}

type ThumbInformation struct {
	ThumbURL    string `json:"thumb_url"`
	ThumbWidth  int    `json:"thumb_width"`
	ThumbHeight int    `json:"thumb_height"`
}

type BaseInlineQueryResults struct {
	ID                  string                `json:"id"`
	Type                string                `json:"type"`
	InputMessageContent InputMessageContent   `json:"input_message_content"`
	ReplyKeyboard       *InlineKeyboardMarkup `json:"reply_keyboard"`
}

type InlineQueryResultArticle struct {
	ThumbInformation
	BaseInlineQueryResults
	Title       string `json:"title"`
	URL         string `json:"url"`
	HideURL     bool   `json:"hide_url"`
	Description string `json:"description"`
}

type InlineQueryResultPhoto struct {
	BaseInlineQueryResults
	PhotoURL        string          `json:"photo_url"`
	ThumbURL        string          `json:"thumb_url"`
	PhotoWidth      int             `json:"photo_width"`
	PhotoHeight     int             `json:"photo_height"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
}

type InlineQueryResultGif struct {
	BaseInlineQueryResults
	GifURL          string          `json:"git_url"`
	GifWidth        int             `json:"gif_width"`
	GifHeight       int             `json:"gif_height"`
	GifDuration     int             `json:"gif_duration"`
	ThumbURL        string          `json:"thumb_url"`
	ThumbMimeType   string          `json:"thumb_mime_type"`
	Title           string          `json:"title"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entity"`
}

type InlineQueryResultMpeg4Gif struct {
	BaseInlineQueryResults
	Mpeg4URL        string          `json:"mpeg4_url"`
	Mpeg4Width      int             `json:"mpeg4_width"`
	Mpeg4Height     int             `json:"mpeg4_height"`
	Mpeg4Duration   int             `json:"mpeg4_duration"`
	ThumbURL        string          `json:"thumb_url"`
	ThumbMimeType   string          `json:"thumb_mime_type"`
	Title           string          `json:"title"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
}

type InlineQueryResultVideo struct {
	BaseInlineQueryResults
	VideoURL        string          `json:"video_url"`
	MimeType        string          `json:"mime_type"`
	ThumbURL        string          `json:"thumb_url"`
	Title           string          `json:"title"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
	VideoWidth      int             `json:"video_width"`
	VideoHeight     int             `json:"video_hieght"`
	VideoDuration   int             `json:"video_duration"`
	Description     string          `json:"descirption"`
}

type InlineQueryResultAudio struct {
	BaseInlineQueryResults
	AudioURL        string          `json:"audio_url"`
	Title           string          `json:"title"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
	Performer       string          `json:"performer"`
	AudioDuration   int             `json:"audio_duration"`
}

type InlineQueryResultVoice struct {
	BaseInlineQueryResults
	VoiceURL        string          `json:"voice_url"`
	Title           string          `json:"title"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
	VoiceDuration   int             `json:"voice_duration"`
}

type InlineQueryResultDocument struct {
	BaseInlineQueryResults
	Title           string          `json:"title"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
	DocumentURL     string          `json:"document_url"`
	MimeType        string          `json:"mime_type"`
	Description     string          `json:"description"`
	ThumbInformation
}

type InlineQueryResultLocation struct {
	BaseInlineQueryResults
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	Title                string  `json:"title"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy"`
	LivePeriod           int     `json:"live_period"`
	Heading              int     `json:"heading"`
	ProximityAlertRadius int     `json:"proximity_alert_radius"`
	ThumbInformation
}

type InlineQueryResultVenue struct {
	BaseInlineQueryResults
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Title           string  `json:"title"`
	Address         string  `json:"address"`
	FoursQuareId    string  `json:"foursquare_id"`
	FoursQuareType  string  `json:"foursquare_type"`
	GooglePlaceID   string  `json:"google_place_id"`
	GooglePlaceType string  `json:"google_place_type"`
	ThumbInformation
}

type InlineQueryResultContact struct {
	BaseInlineQueryResults
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Vcard       string `json:"vcard"`
	ThumbInformation
}

type InlineQueryResultGame struct {
	ID            string                `json:"id"`
	Type          string                `json:"type"`
	GameShortName string                `json:"game_short_name"`
	ReplyKeyboard *InlineKeyboardMarkup `json:"reply_keyboard"`
}

type InlineQueryResultCachedPhoto struct {
	BaseInlineQueryResults
	PhotoFileID     string          `json:"photo_file_id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
}

type InlineQueryResultCachedGif struct {
	BaseInlineQueryResults
	GifFileID       string          `json:"gif_file_id"`
	Title           string          `json:"title"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
}

type InlineQueryResultCachedMpeg4Gif struct {
	BaseInlineQueryResults
	Mpeg4FileID     string          `json:"mpeg4_file_id"`
	Title           string          `json:"title"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
}

type InlineQueryResultCachedSticker struct {
	BaseInlineQueryResults
	StickerFileID   string          `json:"sticker_file_id"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
}

type InlineQueryResultCachedDocument struct {
	DocumentFileID  string          `json:"document_file_id"`
	Description     string          `json:"description"`
	Title           string          `json:"title"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
}

type InlineQueryResultCachedVideo struct {
	VideoFileID     string          `json:"video_file_id"`
	Description     string          `json:"description"`
	Title           string          `json:"title"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
}

type InlineQueryResultCachedVoice struct {
	VoiceFileID     string          `json:"voice_file_id"`
	Title           string          `json:"title"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
}

type InlineQueryResultCachedAudio struct {
	AudioFileID     string          `json:"audio_file_id"`
	Caption         string          `json:"caption"`
	ParseMode       string          `json:"parse_mode"`
	CaptionEntities []MessageEntity `json:"caption_entities"`
}

func (InlineQueryResultCachedMpeg4Gif) getQueryResult() {}
func (InlineQueryResultCachedGif) getQueryResult()      {}
func (InlineQueryResultCachedPhoto) getQueryResult()    {}
func (InlineQueryResultGame) getQueryResult()           {}
func (InlineQueryResultContact) getQueryResult()        {}
func (InlineQueryResultVoice) getQueryResult()          {}
func (InlineQueryResultVenue) getQueryResult()          {}
func (InlineQueryResultLocation) getQueryResult()       {}
func (InlineQueryResultDocument) getQueryResult()       {}
func (InlineQueryResultAudio) getQueryResult()          {}
func (InlineQueryResultVideo) getQueryResult()          {}
func (InlineQueryResultMpeg4Gif) getQueryResult()       {}
func (InlineQueryResultGif) getQueryResult()            {}
func (InlineQueryResultPhoto) getQueryResult()          {}
func (InlineQueryResultArticle) getQueryResult()        {}
