package types

// Sticker ...
type Sticker struct {
	// Width of sticker 512 default
	Width int `json:"width"`

	// Height of sticker, in my case is 502
	Height int `json:"height"`

	// Emoji, idk, is it raises error, or something
	Emoji string `json:"emoji"`

	// SetName is Sticker Pack emoji from
	SetName string `json:"set_name"`

	// IsAnimated ...
	IsAnimated bool `json:"is_animated"`

	// Thumb, idk what is it
	Thumb *Thumb `json:"thumb"`

	// FileID Emoji file id, strange but file_id is string
	FileID string `json:"file_id"`

	// FileUnqiueID you know, just random chars
	FileUnqiueID string `json:"file_unique_id"`

	// Filesize emoji file size in telegram servers
	// showed in bytes, or bites, idk
	FileSize int `json:"file_size"`
}

// Thumb ...
type Thumb struct {
	FileID       string `json:"file_id"`
	FileUnqiueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}
