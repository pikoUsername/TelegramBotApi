package utils

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	// cant make constant
	httpRegex, _    = regexp.Compile("^(http|https)://")
	HTML_QUOTES_MAP = map[interface{}]string{"<": "&lt;", ">": "&gt;", "&": "&amp;", '"': "&quot;"}
)

// Link check out the link for http and https starting with
func Link(link string, text string) (string, error) {
	if !httpRegex.MatchString(link) {
		return "", errors.New("Link is not valid.")
	}
	return fmt.Sprintf("<%s href='%s'>%s</%s>", "a", link, text, "a"), nil
}

// Strong make stronger any text
func Strong(text ...string) string {
	return "<strong>" + fmt.Sprintln(text) + "</strong>"
}

// Italic, spahetti
func HItalic(text ...string) string {
	return "<i>" + fmt.Sprintln(text) + "</i>"
}

// Code is Code, telegram only lanuage- startswith classes for code
func Code(code string, language string) string {
	return fmt.Sprintf("<%s class='language-%s'>%s</%s>", "code", language, code, "code")
}

// Pre pre pre pre pre
func Pre(text ...string) string {
	return "<pre>" + fmt.Sprintln(text) + "</pre>"
}

func PreCode(code string, language string) string {
	return Pre(Code(code, language))
}

func Bold(text ...string) string {
	return "<b>" + fmt.Sprintln(text) + "</b>"
}
