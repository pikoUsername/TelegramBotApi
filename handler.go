package tgp

import (
	"sync"

	"github.com/pikoUsername/tgp/objects"
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
	Handler HandlerFunc
	Filters []interface{}

	mu sync.Mutex
}

func (h *HandlerType) AddFilters(filters ...interface{}) {
	if len(filters) != 0 {
		h.mu.Lock()
		h.Filters = append(h.Filters, filters...)
		h.mu.Unlock()
	}
}

func (h *HandlerType) apply(c *Context) {
	if len(h.Filters) == 0 || checkFilters(h.Filters, c.Update) {
		h.Handler(c)
	}
}

// NewHandlerType returns a HandlerType instance
func NewHandlerType(handler HandlerFunc) *HandlerType {
	return &HandlerType{
		Handler: handler,
		Filters: []interface{}{},
		mu:      sync.Mutex{},
	}
}

// HandlerObj ...
type HandlerObj struct {
	handlers []*HandlerType
	mu       *sync.Mutex
}

// NewHandlerObj creates new DefaultHandlerObj
func NewHandlerObj() *HandlerObj {
	return &HandlerObj{mu: &sync.Mutex{}}
}

func (ho *HandlerObj) Trigger(c *Context) {
	c.handlers = ho.handlers
	for i, h := range ho.handlers {
		if len(h.Filters) == 0 || checkFilters(h.Filters, c.Update) {
			h.Handler(c)
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
		ho.mu.Lock()

		// := using this, type checking will not work
		switch conv := elem.(type) {
		case func(*objects.Update) error:
			partial = true
			handler.AddFilters(conv)

		case func(*Context):
			partial = true
			handler.Handler = conv

		case *HandlerType:
			ho.handlers = append(ho.handlers, conv)
		default:
			return tgpErr.New("only func(*objects.Update), func(*tgp.Context), or *tgp.HandlerType types.")
		}

		ho.mu.Unlock()
	}
	if partial {
		ho.handlers = append(ho.handlers, handler)
	}
	return nil
}
