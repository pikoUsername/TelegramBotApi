package dispatcher

import "github.com/pikoUsername/tgp/bot"

type HandlerObj interface {
	Register(HandlerType)
	Trigger(interface{}, bot.Bot)
}

type HandlerType func(interface{}, bot.Bot)

// HandlerObj uses for save Callback
type DefaultHandlerObj struct {
	Callbacks  []HandlerType
	Middleware Middleware
}

// Register, append to Callbacks, e.g handler functions
func (ho *DefaultHandlerObj) Register(f HandlerType) {
	ho.Callbacks = append(ho.Callbacks, f)
}

// RegisterMiddleware looks like a bad code
// Register middlewares, except function which should return handler
// e.g second parametr
// for example, you want to register every user which writed to you bot
// You can registerMiddleware for MessageHandler, not for all handlers
// Or maybe want to make throttling middleware, just Registers middleware
func (ho *DefaultHandlerObj) RegisterMiddleware(f MiddlewareType) {
	ho.Middleware.Register(f)
}

// Trigger is from aiogram framework
func (ho *DefaultHandlerObj) Trigger(obj interface{}, bot bot.Bot) {
	if ho.Middleware != nil {
		for _, f := range ho.Middleware.GetCallbacks() {
			f(obj, nil) // stub
		}
	}
	for _, cb := range ho.Callbacks {
		cb(obj, bot)
	}
}
