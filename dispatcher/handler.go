package dispatcher

import "github.com/pikoUsername/tgp/objects"

type HandlerObj interface {
	Register(HandlerFunc)
	Trigger(objects.Update)
	RegisterMiddleware(MiddlewareFunc)
}

type HandlerFunc func(objects.Update)

// HandlerObj uses for save Callback
type DefaultHandlerObj struct {
	callbacks  []HandlerFunc
	Middleware MiddlewareManager
}

func NewDHandlerObj(dp *Dispatcher) *DefaultHandlerObj {
	return &DefaultHandlerObj{
		Middleware: NewDMiddlewareManager(dp),
	}
}

// Register, append to Callbacks, e.g handler functions
func (ho *DefaultHandlerObj) Register(f HandlerFunc) {
	ho.callbacks = append(ho.callbacks, f)
}

// RegisterMiddleware looks like a bad code
// for example, you want to register every user which writed to you bot
// You can registerMiddleware for MessageHandler, not for all handlers
// Or maybe want to make throttling middleware, just Registers middleware
//
// Example of middlware see in handler_test.go
func (ho *DefaultHandlerObj) RegisterMiddleware(f MiddlewareFunc) {
	ho.Middleware.Register(f)
}

// Trigger is from aiogram framework
// Trigger is triggers all callbacks in handler
// when middlewares activates, middleware calls a handler
// Just triggers one, you must call Handler in Middleware,
// and handle error, which raised by Handler
func (ho *DefaultHandlerObj) Trigger(upd objects.Update) {
	for _, cb := range ho.callbacks {
		ho.Middleware.Trigger(&upd, cb)
	}
}
