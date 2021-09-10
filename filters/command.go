package filters

import (
	"strings"

	"github.com/pikoUsername/tgp/objects"
)

// Command filter, check out for prefix
// Prefix by default is /, and prefix could be only one character
type CommandFilter struct {
	prefix         string
	cmds           []string
	ignore_mention bool
	ignore_caption bool
}

func (c *CommandFilter) Check(u *objects.Update) bool {
	var mention string
	text := u.Message.Text
	if text == "" {
		return false
	}
	text_args := strings.Split(u.Message.Text, " ")
	raw_text := text_args[0]

	command := strings.ToLower(raw_text)
	prefix := raw_text[:len(c.prefix)]
	if !c.ignore_mention && mention != "" || prefix != c.prefix {
		return false
	}
	command = command[1:]

	for _, cmd := range c.cmds {
		if cmd == command {
			return true
		}
	}
	return false
}

// NewCommand creates new Command object
func Command(cmd ...string) *CommandFilter {
	return &CommandFilter{
		prefix:         "/",
		cmds:           cmd,
		ignore_mention: false,
		ignore_caption: true,
	}
}

// ========================
// Command based filters
// ========================

func CommandStart() *CommandFilter {
	return Command("start")
}

func CommandHelp() *CommandFilter {
	return Command("help")
}

func CommandPrivacy() *CommandFilter {
	return Command("privacy")
}

func CommandSettings() *CommandFilter {
	return Command("settings")
}

func CommandCancel() *CommandFilter {
	return Command("cancel")
}
