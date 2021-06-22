package utils

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	httpRegex = getHTTPRegex()

// 	HTML_QUOTES_MAP = map[interface{}]string{"<": "&lt;", ">": "&gt;", "&": "&amp;", '"': "&quot;"}
)

func getHTTPRegex() *regexp.Regexp {
	regex, _ := regexp.Compile("^(http|https)://")
	return regex
}

type HTMLMarkdown struct{}

// Link check out the link for http and https starting with
func (hm *HTMLMarkdown) Link(link string, text string) (string, error) {
	if !httpRegex.MatchString(link) {
		return "", errors.New("link is not valid")
	}
	return fmt.Sprintf("<%s href='%s'>%s</%s>", "a", link, text, "a"), nil
}

// Strong make stronger any text
func (ht *HTMLMarkdown) Strong(text ...string) string {
	return "<strong>" + fmt.Sprintln(text) + "</strong>"
}

// Italic, spahetti
func (ht *HTMLMarkdown) Italic(text ...string) string {
	return "<i>" + fmt.Sprintln(text) + "</i>"
}

// Code is Code, telegram only lanuage- startswith classes for code
func (hm *HTMLMarkdown) Code(code string, language string) string {
	return fmt.Sprintf(
		"<code class='language-%s'>%s</code>", language, code,
	)
}

// Pre pre pre pre pre
func (hm *HTMLMarkdown) Pre(text ...string) string {
	return "<pre>" + fmt.Sprintln(text) + "</pre>"
}

func (hm *HTMLMarkdown) PreCode(code string, language string) string {
	return hm.Pre(hm.Code(code, language))
}

func (hm *HTMLMarkdown) Bold(text ...string) string {
	return "<b>" + fmt.Sprintln(text) + "</b>"
}

func NewHTMLMarkdown() *HTMLMarkdown {
	return &HTMLMarkdown{}
}

type Markdown2 struct{}

func (md *Markdown2) Link(url string, text ...string) string {
	return "[" + fmt.Sprintln(text) + "](" + url + ")"
}

func (md *Markdown2) Pre(text ...string) string {
	return "```\n" + fmt.Sprintln(text) + "\n```"
}

func (md *Markdown2) PreCode(code string, language string) string {
	return "```" + language + "\n" + code + "```"
}

func (md *Markdown2) Code(language string, text ...string) string {
	return "`" + fmt.Sprintln(text) + "`"
}

func (md *Markdown2) UnderLine(text ...string) string {
	return "__\r" + fmt.Sprintln(text) + "__\r"
}

func (md *Markdown2) StrikeThrough(text ...string) string {
	return "~" + fmt.Sprintln(text) + "~"
}

func (md *Markdown2) Italic(text ...string) string {
	return "_\r" + fmt.Sprintln(text) + "_\r"
}

func (md *Markdown2) Bold(text ...string) string {
	return "*" + fmt.Sprintln(text) + "*"
}

func NewMarkdown2() *Markdown2 {
	return &Markdown2{}
}
