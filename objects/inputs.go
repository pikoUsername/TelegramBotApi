package objects

import (
	"io"
	"os"
	"time"
)

// InputFile ...
type InputFile struct {
	Path   string
	Name   string
	URL    string
	Length int
	File   io.Reader
}

// Read ...
// No compress choice
func (f *InputFile) Read(p []byte) (n int, err error) {
	if f.File != nil && f.URL == "" {
		return f.File.Read(p)
	}
	file, err := os.Open(f.Path)
	if err != nil {
		return 0, err
	}
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	// you should call Close after Read function
	f.File = file
	f.Length = (int)(stat.Size())
	return f.Length, nil
}

func (f *InputFile) Close() error {
	if s, ok := f.File.(io.ReadCloser); ok {
		return s.Close()
	}
	return nil
}

func NewInputFile(path, name string) (*InputFile, error) {
	var err error

	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return &InputFile{}, err
	}
	stat, err := f.Stat()
	if err != nil {
		return &InputFile{}, err
	}
	return &InputFile{
		File:   f,
		Path:   path,
		Name:   name,
		Length: int(stat.Size()),
	}, nil
}

func NewInputFileFromReader(r io.Reader, length int, name string) *InputFile {
	return &InputFile{
		File:   r,
		Length: length,
	}
}

type InputMedia interface {
	GetMedia() string
}

type InputMediaVideo struct {
	Type  string `json:"type"`
	Media string `json:"media"`

	// TODO: use it with InputFile
	Thumb             string           `json:"thumb"`
	Caption           string           `json:"caption"`
	ParseMode         string           `json:"parse_mode"`
	CaptionEntities   []*MessageEntity `json:"caption_entities"`
	Width             int64            `json:"width"`
	Height            int64            `json:"height"`
	Duration          time.Duration    `json:"duration"`
	SupportsStreaming bool             `json:"supports_streaming"`
}

type InputMediaAnimation struct {
}

type InputMediaDocument struct {
}

type InputMediaVoice struct {
}
