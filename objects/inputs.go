package objects

import (
	"io"
	"io/ioutil"
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
	bs, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return 0, err
	}
	p = bs
	bs_len := len(bs)
	f.Length = bs_len
	return bs_len, nil
}

func NewInputFile(path string, name string) (*InputFile, error) {
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
		Name:   name,
		Path:   path,
		Length: int(stat.Size()),
	}, nil
}

func InputFileFromReader(r io.Reader, length int, name string) *InputFile {
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
