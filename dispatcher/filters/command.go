package filters

import (
	"strings"

	"github.com/pikoUsername/tgp/objects"
)

// Command filter, check out for prefix
// Prefix by default is /
type Command struct {
	prefix         string
	cmds           []string
	ignore_mention bool
	ignore_caption bool
}

func (c *Command) Check(u *objects.Update) bool {
	text := u.Message.Text
	if text == "" {
		return false
	}
	text_args := strings.Split(u.Message.Text, " ")
	args := text_args[0]
	mention := strings.ToLower(strings.Split(args, "@")[1])

	if !c.ignore_caption && mention != "" {
		return false
	}
	return true
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
