package dispatcher

import (
	"github.com/pikoUsername/tgp/bot"
	"github.com/pikoUsername/tgp/configs"
	"github.com/pikoUsername/tgp/objects"
)

type Context struct {
	bot     *bot.Bot
	Message *objects.Message
}

func NewContext(bot *bot.Bot, mes *objects.Message) *Context {
	return &Context{
		bot:     bot,
		Message: mes,
	}
}

// Send method alias for sending using bot
func (mc *Context) Send(c configs.Configurable) (*objects.Message, error) {
	return mc.bot.Send(c)
}
