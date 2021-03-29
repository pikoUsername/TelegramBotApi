package objects

// File represents telegram File object
// https://core.telegram.org/bots/api#file
type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`

	// Note:
	//  FilePath is Url to file
	FilePath string `json:"file_path"`
}
