package filters

import (
	"strings"

	"github.com/pikoUsername/tgp/objects"
)

// Command filter, check out for prefix
// Prefix by default is /, and prefix could be only one character
type Command struct {
	prefix         string
	cmds           []string
	ignore_mention bool
	ignore_caption bool
}

func (c *Command) Check(u *objects.Update) bool {
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
func NewCommand(cmd ...string) *Command {
	return &Command{
		prefix:         "/",
		cmds:           cmd,
		ignore_mention: false,
		ignore_caption: true,
	}
}

// ========================
// Command based filters
// ========================

func NewCommandStart() *Command {
	return NewCommand("start")
}

func NewCommandHelp() *Command {
	return NewCommand("help")
}

func NewCommandPrivacy() *Command {
	return NewCommand("privacy")
}

func NewCommandSettings() *Command {
	return NewCommand("settings")
}

func NewCommandCancel() *Command {
	return NewCommand("cancel")
}
