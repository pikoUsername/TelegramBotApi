package dispatcher

import "github.com/pikoUsername/tgp/bot"

// HandlerObj uses for save Callback
type HandlerObj struct {
	Callbacks   []*func(interface{}, bot.Bot)
	Middlewares []func(interface{}, *func(interface{}, bot.Bot)) *func(interface{}, bot.Bot)
}

// Register, append to Callbacks, e.g handler functions
func (ho *HandlerObj) Register(f *func(interface{}, bot.Bot)) {
	ho.Callbacks = append(ho.Callbacks, f)
}

// RegisterMiddleware looks like a bad code
// Register middlewares, except function which should return handler
// e.g second parametr
// for example, you want to register every user which writed to you bot
// You can registerMiddleware for MessageHandler, not for all handlers
// Or maybe want to make throttling middleware, just Registers middleware
func (ho *HandlerObj) RegisterMiddleware(f func(interface{}, *func(interface{}, bot.Bot)) *func(interface{}, bot.Bot)) {
	ho.Middlewares = append(ho.Middlewares, f)
}

// Trigger is from aiogram framework
func (ho *HandlerObj) Trigger(obj interface{}, bot bot.Bot) {
	if ho.Middlewares != nil {
		for _, f := range ho.Middlewares {
			f(obj, nil) // stub
		}
	}

}
