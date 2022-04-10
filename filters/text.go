package filters

import (
	"strings"

	"github.com/pikoUsername/tgp/objects"
)

type TextFilter struct {
	Text        string
	Ignore_case bool
	Equals      bool
	Contains    bool
	Endswith    bool
	Startswith  bool
}

func (t *TextFilter) Check(u *objects.Update) bool {
	var text string

	// CheckOut for text
	if u.Message != nil {
		text = u.Message.Text
	} else if u.CallbackQuery != nil {
		text = u.CallbackQuery.Data
	} else if u.InlineQuery != nil {
		text = u.InlineQuery.Query
	} else if u.Poll != nil {
		text = u.Poll.Question
	}

	if t.Ignore_case {
		text = strings.ToLower(text)
	}

	if t.Equals {
		return t.Text == text
	} else if t.Contains {
		return strings.Contains(text, t.Text)
	} else if t.Endswith {
		return strings.HasSuffix(text, t.Text)
	} else if t.Startswith {
		return strings.HasPrefix(text, t.Text)
	}

	return false
}

func Text(text string) *TextFilter {
	return &TextFilter{
		Text:        text,
		Ignore_case: true,
		// other fields is false, by default, boolean type by default is false
	}
}
