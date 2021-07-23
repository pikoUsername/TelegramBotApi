package objects

// ForceReply ...
type ForceReply struct {
	ForceReply            bool   `json:"force_reply"`
	InputFieldPlaceHolder string `json:"input_field_placeholder"`
	Selective             bool   `json:"selective"`
}
