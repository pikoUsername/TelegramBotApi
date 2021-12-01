package tgp

import (
	"sync"
)

type HandlerFunc func(*Context)

// errors
var (
	MiddlewareTypeInvalid = tgpErr.New("typ parameter of variable not in ['post', 'pre', 'process']")
	MiddlewareNotFound    = tgpErr.New("passed middleware not found")
	MiddlewareIncorrect   = tgpErr.New("passed function is not function type")
)

// Another level of abstraction
// Filters field is interface{}, types:
// func(u *objects.Update) and Filter interface
//
// ```go
// 	dp.MessageHandler.Register(
//		func(ctx *tgp.Context) {...}, // handler
// 		func(upd *objects.Update) {return u.Message.From.ID == <owner_id>}, // filter
// 	)
// ```
// Can be used as handlers chain, e.g registered by ResgisterChain()
type HandlerType struct {
	handler HandlerFunc
	filters []interface{}

	mu sync.Mutex
}

func (h *HandlerType) AddFilters(filters ...interface{}) {
	if len(filters) != 0 {
		h.mu.Lock()
		h.filters = append(h.filters, filters...)
		h.mu.Unlock()
	}
}

func (ht *HandlerType) SetHandler(h HandlerFunc) {
	ht.handler = h
}

func (h *HandlerType) apply(c *Context) {
	if len(h.filters) == 0 || checkFilters(h.filters, c.Update) {
		h.handler(c)
	}
}

// NewHandlerType returns a HandlerType instance
func NewHandlerType(handler HandlerFunc) *HandlerType {
	return &HandlerType{
		handler: handler,
		filters: []interface{}{},
		mu:      sync.Mutex{},
	}
}

// HandlerObj ...
type HandlerObj struct {
	handlers []*HandlerType
	mu       sync.Mutex
}

// NewHandlerObj creates new DefaultHandlerObj
func NewHandlerObj() *HandlerObj {
	return &HandlerObj{mu: sync.Mutex{}}
}

func (ho *HandlerObj) Trigger(c *Context) {
	c.handlers = ho.handlers
	for i, h := range ho.handlers {
		if len(h.filters) == 0 || checkFilters(h.filters, c.Update) {
			h.handler(c)
			c.cursor = i
		}
	}
}

// Register could work in two modes
//
// 1. registers a filters, and handlers, this mode will register chain in the end of function
// 2. registers handlerchain instantly
// works for every functions argument
func (ho *HandlerObj) Register(callbacks ...interface{}) error {
	var partial bool
	handler := &HandlerType{}

	for _, elem := range callbacks {
		switch conv := elem.(type) {
		case Filter:
			partial = true
			handler.AddFilters(conv)

		case func(*Context):
			partial = true
			handler.handler = conv

		case *HandlerType:
			ho.mu.Lock()
			ho.handlers = append(ho.handlers, conv)
			ho.mu.Unlock()
		default:
			return tgpErr.New("only func(*objects.Update), func(*tgp.Context), or *tgp.HandlerType types.")
		}
	}
	if partial {
		ho.mu.Lock()
		ho.handlers = append(ho.handlers, handler)
		ho.mu.Unlock()
	}
	return nil
}
