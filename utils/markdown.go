package utils

import "fmt"

// You will use is most cases HTML formatting
// So, not need to make MarkDown too? yes?
type HTMLarguments map[string]string

func (hargs *HTMLarguments) Add(key string, value interface{}) {
	return
}

func (hargs *HTMLarguments) Encode() string {
	return ""
}

// wrapHtml wrapps into HTML tag, and tag arguments
func wrapHtml(tag string, text string, hargs *HTMLarguments) string {
	if hargs == nil {
		// wrapps into html tag
		return fmt.Sprintf("<%s>%s</%s>", tag, text, tag)
	} else {
		return fmt.Sprintf("<%s %s>%s</%s>", tag, hargs.Encode(), text, tag)
	}
}

func link(link string, text string) string {
	return wrapHtml("a", text, &HTMLarguments{})
}
