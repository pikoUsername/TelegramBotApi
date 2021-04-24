package objects

type Video struct {
	BaseFile
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
}
