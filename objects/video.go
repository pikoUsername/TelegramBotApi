package objects

import "time"

type Video struct {
	BaseFile
	Width    int           `json:"width"`
	Height   int           `json:"height"`
	Duration time.Duration `json:"duration"`
	MimeType string        `json:"mime_type"`
}
