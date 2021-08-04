package tgp

import (
	"reflect"

	"github.com/pikoUsername/tgp/objects"
)

// HandlerFunc is stub
type HandlerFunc interface{}

// Another level of abstraction
// Filters field is interface{}, types:
// func(u *objects.Update) and Just Filter interface
//
// example: ```go
// 	dp.MessageHandler.Register(
//		func(m *objects.Message) {...},
// 		func(u *objects.Update) {return u.Message.From.ID == <owner_id>},
// 	)
// ```
type HandlerType struct {
	Callback HandlerFunc
	Filters  []interface{}
}

// CheckForFilters iterate all filters and call Check method for check
func (ht *HandlerType) CheckForFilters(u *objects.Update) bool {
	var b bool

	for _, f := range ht.Filters {
		switch t := f.(type) {
		case Filter:
			b = t.Check(u)
		case func(u *objects.Update) bool:
			b = t(u)
		default:
			continue
		}
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

// HandlerObj uses for save Callback
type HandlerObj struct {
	handlers   []*HandlerType
	Middleware MiddlewareManager
}

// NewHandlerObj creates new DefaultHandlerObj
func NewHandlerObj(dp *Dispatcher) *HandlerObj {
	return &HandlerObj{
		Middleware: NewMiddlewareManager(dp),
	}
}

// Register, append to Callbacks, e.g handler functions
func (ho *HandlerObj) Register(f HandlerFunc, filters ...interface{}) {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return
	}

	ht := HandlerType{
		Callback: &f,
		Filters:  filters,
	}

	ho.handlers = append(ho.handlers, &ht)
}

// Unregister checkout to memory address
// and cut up it if find something, with same address
func (ho *HandlerObj) Unregister(handler *HandlerFunc) {
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
func (ho *HandlerObj) RegisterMiddleware(f ...MiddlewareFunc) {
	ho.Middleware.Register(f...)
}

func (ho *HandlerObj) TriggerMiddleware(bot *Bot, update *objects.Update, typ string) error {
	return ho.Middleware.Trigger(bot, update, typ)
}
