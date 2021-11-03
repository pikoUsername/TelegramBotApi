package filters

import (
	"regexp"

	"github.com/pikoUsername/tgp/objects"
)

type RegexpFilter struct {
	Regexp *regexp.Regexp
}

func (r *RegexpFilter) Check(u *objects.Update) bool {
	var content string
	if u.Message != nil {
		content = u.Message.Text
	} else if u.CallbackQuery != nil {
		content = u.CallbackQuery.Message.Text
	} else if u.Poll != nil {
		content = u.Poll.Question
	} else {
		return false
	}

	match := string(r.Regexp.Find([]byte(content)))
	return match != ""
}

func Regexp(re string) (*RegexpFilter, error) {
	rex, err := regexp.Compile(re)
	if err != nil {
		return &RegexpFilter{}, err
	}

	return &RegexpFilter{
		Regexp: rex,
	}, nil
}
