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
// Middleware can be used as error handlers
// There are 2 middleware types.
// First is middleware that wrapps a handler.
// Second one is just handler that uses Next method to go next handler
// (i m not sure, can this be considered as middleware and not fsm?)
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

type HandlerChain interface {
	Trigger(*Context)
	HandlerFunc(HandlerFunc) *HandlerType
	Handlers() []HandlerFunc
	Use(md ...MiddlewareFunc)
}

type DefaultHandlerChain struct {
	middleware []MiddlewareFunc
	handlers   []*HandlerType
	mu         sync.Mutex
}

// NewHandlerChain creates new DefaultHandlerChain
func NewHandlerChain() *DefaultHandlerChain {
	return &DefaultHandlerChain{mu: sync.Mutex{}}
}

func (ho *DefaultHandlerChain) Handlers() []HandlerFunc {
	// copies handlers object
	l := make([]HandlerFunc, len(ho.handlers))
	for _, h := range ho.handlers {
		l = append(l, h.handler)
	}
	return l
}

func (ho *DefaultHandlerChain) Trigger(c *Context) {
	c.handlers = ho.handlers
	for i, h := range ho.handlers {
		if len(h.filters) == 0 || checkFilters(h.filters, c.Update) {
			if len(ho.middleware) > 0 {
				for _, md := range ho.middleware {
					md(h.handler)(c)
				}
			} else {
				h.handler(c)
			}
			c.cursor = i
			break
		}
	}
}

func (ho *DefaultHandlerChain) Use(md ...MiddlewareFunc) {
	ho.middleware = append(ho.middleware, md...)
}

// HandlerFunc appends new handlerType, and returns it
func (ho *DefaultHandlerChain) HandlerFunc(h HandlerFunc) *HandlerType {
	handler := &HandlerType{handler: h}
	ho.handlers = append(ho.handlers, handler)
	return handler
}
