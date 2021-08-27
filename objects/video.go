package objects

import "time"

type Video struct {
	*BaseFile
	MimeType string        `json:"mime_type"`
	Width    int           `json:"width"`
	Height   int           `json:"height"`
	Duration time.Duration `json:"duration"`
}
