package objects

import "time"

// Voice represents Voice telegram object
// https://core.telegram.org/bots/api#voice
type Voice struct {
	FileID       string        `json:"file_id"`
	FileUniqueID string        `json:"file_unique_id"`
	Duration     time.Duration `json:"duration"`
	MimeType     string        `json:"mime_type"`
	FileSize     int64         `json:"file_size"`
}
