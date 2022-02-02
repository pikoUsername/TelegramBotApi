package tgp

import (
	"sync"

	"github.com/pikoUsername/tgp/filters"
)

type HandlerFunc func(*Context)

// Another level of abstraction
// Filters field is Filter interface
// func(u *objects.Update) and Filter interface
//
// ```
// dp.MessageHandler.HandleFunc(func(ctx *tgp.Context) {
// 	   ctx.Reply("HERE")
// }).Command("/start")
// ```
// Can be used as handlers chain, e.g registered by ResgisterChain()
type HandlerType struct {
	handler HandlerFunc
	filters []Filter
}

func (he *HandlerType) Filters(filters ...Filter) *HandlerType {
	he.filters = append(he.filters, filters...)
	return he
}

// Regexp register regexp filter
// note: errors will be ingored
func (he *HandlerType) Regexp(pattern string) *HandlerType {
	filter, _ := filters.Regexp(pattern)
	he.filters = append(he.filters, filter)
	return he
}

func (he *HandlerType) Text(pattern string) *HandlerType {
	filter := filters.Text(pattern)
	he.filters = append(he.filters, filter)
	return he
}

func (he *HandlerType) Command(cmd string) *HandlerType {
	filter := filters.Command(cmd)
	he.filters = append(he.filters, filter)
	return he
}

// GetFilters returns filters interfaces, types tgp.FilterFunc, tgp.Filter
func (he *HandlerType) GetFilters() []Filter {
	return he.filters
}

func (he *HandlerType) GetHandler() HandlerFunc {
	return he.handler
}

// Use middleware registration implementation
func (he *HandlerType) Use(middleware MiddlewareFunc) *HandlerType {
	he.handler = middleware(he.handler)
	return he
}

func (he *HandlerType) Handler(handler HandlerFunc) *HandlerType {
	he.handler = handler
	return he
}

func (he *HandlerType) Copy() *HandlerType {
	ht := &HandlerType{}
	if he.handler != nil {
		ht.handler = he.handler
	}
	copy(ht.filters, he.filters)
	return ht
}

// NewHandlerType returns a HandlerType instance
func NewHandlerType(handler HandlerFunc) *HandlerType {
	return &HandlerType{
		handler: handler,
		filters: []Filter{},
	}
}

type HandlerObj interface {
	Trigger(*Context)
	HandlerFunc(HandlerFunc) *HandlerType
	Register(...interface{}) *HandlerType
	Handlers() []HandlerFunc
}

type DefaultHandlerObj struct {
	handlers []*HandlerType
	mu       sync.Mutex
}

// NewHandlerObj creates new DefaultHandlerObj
func NewHandlerObj() *DefaultHandlerObj {
	return &DefaultHandlerObj{mu: sync.Mutex{}}
}

func (ho *DefaultHandlerObj) Handlers() []HandlerFunc {
	// copies handlers object
	l := make([]HandlerFunc, len(ho.handlers))
	for _, h := range ho.handlers {
		l = append(l, h.handler)
	}
	return l
}

func (ho *DefaultHandlerObj) Trigger(c *Context) {
	c.handlers = ho.handlers
	for i, h := range ho.handlers {
		if len(h.filters) == 0 || checkFilters(h.filters, c.Update) {
			h.handler(c)
			c.cursor = i
			break
		}
	}
}

// HandlerFunc appends new handlerType, and returns it
func (ho *DefaultHandlerObj) HandlerFunc(h HandlerFunc) *HandlerType {
	handler := &HandlerType{handler: h}
	ho.handlers = append(ho.handlers, handler)
	return handler
}

// Register could work in two modes
//
// 1. registers a filters, and handlers, this mode will register chain in the end of function
// 2. registers handlerchain instantly
// works for every functions argument
func (ho *DefaultHandlerObj) Register(callbacks ...interface{}) *HandlerType {
	var partial bool
	handler := &HandlerType{}

	for _, elem := range callbacks {
		switch conv := elem.(type) {
		case Filter:
			partial = true
			handler.Filters(conv)

		case func(*Context):
			partial = true
			handler.Handler(conv)

		case *HandlerType:
			ho.mu.Lock()
			ho.handlers = append(ho.handlers, conv)
			ho.mu.Unlock()
		default:
			return nil
		}
	}
	if partial {
		ho.mu.Lock()
		ho.handlers = append(ho.handlers, handler)
		ho.mu.Unlock()
	}
	return handler
}
