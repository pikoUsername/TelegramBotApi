package objects

type MaskPosition struct {
	Point string `json:"point"`
	// i dont sure is it safe?
	XShift float32 `json:"x_shift"`
	YShift float32 `json:"y_shift"`
	Scale  float32 `json:"scale"`
}

// Sticker ...
type Sticker struct {
	FileId int64 `json:"file_id"`

	FileUniqueID int64 `json:"file_unique_id"`

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

	IsVideo bool `json:"is_video"`

	// Thumb, idk what is it
	Thumb *Thumb `json:"thumb"`

	// FileID Emoji file id, strange but file_id is string
	FileID string `json:"file_id"`

	// FileUnqiueID you know, just random chars
	FileUnqiueID string `json:"file_unique_id"`

	MaskPosition MaskPosition `json:"mask_position"`

	// Filesize emoji file size in telegram servers
	// showed in bytes, or bites, idk
	FileSize int `json:"file_size"`
}

type StickerSet struct {
	Name          string     `json:"name"`
	Title         string     `json:"title"`
	IsAnimated    bool       `json:"is_animated"`
	IsVideo       bool       `json:"is_video"`
	ContainsMasks bool       `json:"contains_masks"`
	Stickers      []*Sticker `json:"stickers"`
	Thumb         PhotoSize  `json:"thumb"`
}

// Thumb ...
type Thumb struct {
	FileID       string `json:"file_id"`
	FileUnqiueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}
