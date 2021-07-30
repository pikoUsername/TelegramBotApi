package tgp

// ReadFromInputFile call InputFile.Read method and for more shorter code this function was created
// using in Values() in SendDocumentConfig struct
func ReadFromInputFile(v *InputFile, compress bool) (p []byte, err error) {
	bs := make([]byte, v.Length)
	_, err = v.Read(bs)
	return bs, err
}
