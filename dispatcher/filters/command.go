package filters

import (
	"strings"

	"github.com/pikoUsername/tgp/objects"
)

type Command struct {
	cmd string
}

func (c *Command) Check(u *objects.Update) bool {
	return strings.HasPrefix(u.Message.Text, c.cmd)
}
