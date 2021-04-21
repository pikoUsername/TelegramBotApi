package dispatcher

import (
	"github.com/pikoUsername/tgp/bot"
)

type HandlerObj interface {
	Register(HandlerType)
	Trigger(interface{}, *bot.Bot) error
}

type HandlerType func(interface{}, bot.Bot) error

// HandlerObj uses for save Callback
type DefaultHandlerObj struct {
	callbacks  []HandlerType
	Middleware MiddlewareManager
}

func NewDHandlerObj(dp *Dispatcher) *DefaultHandlerObj {
	return &DefaultHandlerObj{
		Middleware: NewDMiddlewareManager(dp),
	}
}

// Register, append to Callbacks, e.g handler functions
func (ho *DefaultHandlerObj) Register(f HandlerType) {
	ho.callbacks = append(ho.callbacks, f)
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
// Trigger is triggers all callbacks in handler
// when middlewares activates, middleware calls a handler
func (ho *DefaultHandlerObj) Trigger(obj interface{}, bot *bot.Bot) error {
	for _, cb := range ho.callbacks {
		err := cb(obj, *bot)
		ho.Middleware.Trigger(obj, &cb)
		if err != nil {
			return err
		}
	}
	return nil
}
