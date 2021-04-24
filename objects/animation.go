package objects

// Animation represents Animation object
// https://core.telegram.org/bots/api#animation
type Animation struct {
	BaseFile
	Width    int `json:"width"`
	Height   int `json:"height"`
	Duration int `json:"duration"`
	*Thumb   `json:"thumb"`
	MimeType string `json:"mime_type"`
}
