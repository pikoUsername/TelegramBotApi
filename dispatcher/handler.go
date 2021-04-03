package dispatcher

import "github.com/pikoUsername/tgp/bot"

// HandlerObj uses for save Callback
type HandlerObj struct {
	Callbacks   []*func(interface{}, *bot.Bot)
	Middlewares []*func(interface{})
}

func (ho *HandlerObj) Register(f *func(interface{}, *bot.Bot)) {
	ho.Callbacks = append(ho.Callbacks, f)
}

func (ho *HandlerObj) RegisterMiddleware(f *func(interface{})) {
	ho.Middlewares = append(ho.Middlewares, f)
}

// Trigger is from aiogram framework
func (ho *HandlerObj) Trigger(obj interface{}) {
	return
}
