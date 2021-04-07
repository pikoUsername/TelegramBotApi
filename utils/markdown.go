package utils

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	httpRegex, _ = regexp.Compile("^(http|https)://")
)

// You will use is most cases HTML formatting
// So, not need to make MarkDown too? yes?
type HTMLTagarguments map[string]string

func (hargs *HTMLTagarguments) Add(key string, value interface{}) {
	return
}

func (hargs *HTMLTagarguments) Encode() string {
	return ""
}

// wrapHtml wrapps into HTML tag, and tag arguments
// Provate method, not for using outside
func wrapHtml(tag string, text string, hargs *HTMLTagarguments) string {
	if hargs == nil {
		// wrapps into html tag
		return fmt.Sprintf("<%s>%s</%s>", tag, text, tag)
	} else {
		return fmt.Sprintf("<%s %s>%s</%s>", tag, hargs.Encode(), text, tag)
	}
}

func Link(link string, text string) (string, error) {
	if !httpRegex.MatchString(link) {
		return text, errors.New("Link is not valid.")
	}
	return wrapHtml("a", text, nil), nil
}

func Strong(text string) string {
	return wrapHtml("strong", text, nil)
}

func HItalic(text string) string {
	return wrapHtml("i", text, nil)
}

func Code(code string, language string) string {
	v := &HTMLTagarguments{}
	v.Add("class", "language-"+language)
	return wrapHtml("code", code, v)
}

func PreCode(code string, language string) string {
	return wrapHtml("pre", Code(code, language), nil)
}

func Pre(text string) string {
	return wrapHtml("pre", text, nil)
}

func Bold(text string) string {
	return wrapHtml("b", text, nil)
}
