package filters

import "github.com/pikoUsername/tgp/objects"

type ContentType struct {
	ctype string
}

func (ct *ContentType) Check(u *objects.Update) bool {
	return u.Message.GetContentType() == ct.ctype
}

func NewContentType(ctype string) *ContentType {
	return &ContentType{
		ctype: ctype,
	}
}
