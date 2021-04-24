package objects

type Document struct {
	BaseFile
	*Thumb   `json:"thumb"`
	MimeType string `json:"mime_type"`
}
