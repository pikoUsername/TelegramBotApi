package objects

import (
	"io"
	"io/ioutil"
	"os"
)

// InputFile ...
// NOTE:
// 	  InputFile.Name must be same as field that stores InputFile.
type InputFile struct {
	Path   string
	URL    string
	Length int
	Name   string
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
		Path:   path,
		Name:   name,
		Length: int(stat.Size()),
	}, nil
}

func InputFileFromReader(r io.Reader, length int, name string) *InputFile {
	return &InputFile{
		File:   r,
		Length: length,
		Name:   name,
	}
}

type InputFilePlaceHolder interface {
	GetMedia() string
}

type InputMediaAudio struct {
	Type  string
	Media string
	Thumb
}

type InputMediaDocument struct {
}

type InputMediaVideo struct {
}
