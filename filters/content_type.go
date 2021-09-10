package filters

import "github.com/pikoUsername/tgp/objects"

type ContentTypeFilter struct {
	ctype string
}

func (ct *ContentTypeFilter) Check(u *objects.Update) bool {
	return u.Message.GetContentType() == ct.ctype
}

func ContentType(ctype string) *ContentTypeFilter {
	return &ContentTypeFilter{
		ctype: ctype,
	}
}
