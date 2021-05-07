package dispatcher

import "github.com/pikoUsername/tgp/objects"

type HandlerFunc func(*objects.Update)

// Another level of abstraction
type HandlerType struct {
	Callback HandlerFunc
	Filters  []Filter
}

// CheckForFilters iterate all filters and call Check method for check
func (ht *HandlerType) CheckForFilters(u *objects.Update) bool {
	for _, f := range ht.Filters {
		b := f.Check(u)
		if !b {
			return false
		}
	}
	return true
}

// Call uses for checking using filters
func (ht *HandlerType) Call(u *objects.Update) {
	fr := ht.CheckForFilters(u)
	if !fr {
		return
	}

	if u != nil {
		ht.Callback(u)
	}
}

// Interface for creating custom HandlerObj
type HandlerObj interface {
	Register(HandlerFunc, ...Filter)
	Trigger(objects.Update)
	RegisterMiddleware(...MiddlewareFunc)
}

// HandlerObj uses for save Callback
type DefaultHandlerObj struct {
	handlers   []HandlerType
	Middleware MiddlewareManager
}

// NEwDHandlerObj creates new DefaultHandlerObj
func NewDHandlerObj(dp *Dispatcher) *DefaultHandlerObj {
	return &DefaultHandlerObj{
		Middleware: NewDMiddlewareManager(dp),
	}
}

// Register, append to Callbacks, e.g handler functions
func (ho *DefaultHandlerObj) Register(f HandlerFunc, filters ...Filter) {
	ht := HandlerType{
		Callback: f,
		Filters:  filters,
	}

	ho.handlers = append(ho.handlers, ht)
}

// RegisterMiddleware looks like a bad code
// for example, you want to register every user which writed to you bot
// You can registerMiddleware for MessageHandler, not for all handlers
// Or maybe want to make throttling middleware, just Registers middleware
//
// Example of middlware see in handler_test.go
func (ho *DefaultHandlerObj) RegisterMiddleware(f ...MiddlewareFunc) {
	ho.Middleware.Register(f...)
}

// Trigger is from aiogram framework
// Trigger is triggers all callbacks in handler
// when middlewares activates, middleware calls a handler
// Just triggers one, you must call Handler in Middleware,
// and handle error, which raised by Handler
func (ho *DefaultHandlerObj) Trigger(upd objects.Update) {
	for _, cb := range ho.handlers {
		ho.Middleware.Trigger(&upd, cb)
	}
}
