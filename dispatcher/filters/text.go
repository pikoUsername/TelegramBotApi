package filters

import (
	"strings"

	"github.com/pikoUsername/tgp/objects"
)

type Text struct {
	text        string
	ignore_case bool
	equals      bool
	contains    bool
	endswith    bool
	startswith  bool
}

func (t *Text) Check(u *objects.Update) bool {
	var text string

	// CheckOut for text
	if u.Message != nil {
		text = u.Message.Text
	} else if u.CallbackQuery != nil {
		// TODO: Callback query object
		// text = u.CallbackQuery
	} else if u.InlineQuery != nil {
		text = u.InlineQuery.Query
	} else if u.Poll != nil {
		text = u.Poll.Question
	}

	if t.ignore_case {
		text = strings.ToLower(text)
	}

	if t.equals {
		return t.text == text
	} else if t.contains {
		return strings.Contains(text, t.text)
	} else if t.endswith {
		return strings.HasSuffix(text, t.text)
	} else if t.startswith {
		return strings.HasPrefix(text, t.text)
	}

	return false
}

func NewText(text string) *Text {
	return &Text{
		text:        text,
		ignore_case: true,
		// other fields is false, by default, boolean type by default is false
	}
}
