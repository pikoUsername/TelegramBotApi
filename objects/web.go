package objects

type WebAppInfo struct {
	URL string `json:"url"`
}

type SentWebAppMessage struct {
	InlineMessageID string `json:"inline_message_id"`
}

type WebAppData struct {
	Data       string `json:"data"`
	ButtonText string `json:"button_text"`
}
