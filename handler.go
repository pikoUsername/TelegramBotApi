package tgp

import (
	"github.com/pikoUsername/tgp/objects"
)

// HandlerFunc is stub
type HandlerFunc interface{}

// Another level of abstraction
type HandlerType struct {
	Callback *HandlerFunc
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
func (ht *HandlerType) Call(u *objects.Update, f func(), sync bool) {
	fr := ht.CheckForFilters(u)
	if !fr {
		return
	}

	if sync {
		f()
	} else {
		go f()
	}
}

// Interface for creating custom HandlerObj
type HandlerObj interface {
	Register(handler HandlerFunc, filters ...Filter)
	Unregister(handler *HandlerFunc)
	RegisterMiddleware(middlewares ...MiddlewareFunc)
	GetHandlers() []HandlerType
	TriggerMiddleware(update *objects.Update, typ string) error
}

// HandlerObj uses for save Callback
type DefaultHandlerObj struct {
	handlers   []HandlerType
	Middleware MiddlewareManager
}

// NEwDHandlerObj creates new DefaultHandlerObj
func NewDHandlerObj(dp *Dispatcher) *DefaultHandlerObj {
	return &DefaultHandlerObj{
		Middleware: NewMiddlewareManager(dp),
	}
}

// Register, append to Callbacks, e.g handler functions
func (ho *DefaultHandlerObj) Register(f HandlerFunc, filters ...Filter) {
	ht := HandlerType{
		Callback: &f,
		Filters:  filters,
	}

	ho.handlers = append(ho.handlers, ht)
}

// Unregister checkout to memory address
// and cut up it if find something, with same address
func (ho *DefaultHandlerObj) Unregister(handler *HandlerFunc) {
	var index int
	for i, h := range ho.handlers {
		if h.Callback == handler {
			// deleting from slice
			index = i - 1
			if index < 0 {
				index = 0
			}
			ho.handlers = append(ho.handlers[:index], ho.handlers[i:]...)
		}
	}
}

// RegisterMiddleware ...
// for example, you want to register every user which writed to you bot
// You can registerMiddleware for MessageHandler, not for all handlers
// Or maybe want to make throttling middleware, just Registers middleware
//
// Example of middlware see in handler_test.go
func (ho *DefaultHandlerObj) RegisterMiddleware(f ...MiddlewareFunc) {
	ho.Middleware.Register(f...)
}

func (ho *DefaultHandlerObj) GetHandlers() []HandlerType {
	return ho.handlers
}

func (ho *DefaultHandlerObj) TriggerMiddleware(update *objects.Update, typ string) error {
	return ho.Middleware.Trigger(update, typ)
}
