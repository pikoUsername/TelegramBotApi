package tgp

import (
	"sync"
	"unsafe"

	"github.com/pikoUsername/tgp/objects"
)

type HandlerFunc func(ctx *Context)

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
	Handlers []HandlerFunc
	Filters  []interface{}

	mu sync.Mutex
}

func (h *HandlerType) Add(handlers ...HandlerFunc) {
	if len(handlers) != 0 {
		h.mu.Lock()
		h.Handlers = append(h.Handlers, handlers...)
		h.mu.Unlock()
	}
}

func (h *HandlerType) AddFilters(filters ...interface{}) {
	if len(filters) != 0 {
		h.mu.Lock()
		h.Filters = append(h.Filters, filters...)
		h.mu.Unlock()
	}
}

// NewHandlerType returns a HandlerType instance
func NewHandlerType(handlers ...HandlerFunc) *HandlerType {
	return &HandlerType{
		Handlers: handlers,
		Filters:  []interface{}{},
		mu:       sync.Mutex{},
	}
}

// HandlerObj uses for save Callback
type HandlerObj struct {
	handlers   []*HandlerType
	mu         *sync.Mutex
	handlerPos uint16
}

// NewHandlerObj creates new DefaultHandlerObj
func NewHandlerObj() *HandlerObj {
	return &HandlerObj{mu: &sync.Mutex{}}
}

func (ho *HandlerObj) Trigger(c *Context) {
	var filtersResult bool
	var nextIndex int
	for i, h := range ho.handlers {
		if h != nil {
			ho.handlerPos = *(*uint16)(unsafe.Pointer(&i))
			nextIndex = i - 1
			if nextIndex < len(ho.handlers) {
				nextIndex = i
			}
			handlers := ho.handlers[nextIndex].Handlers
			for _, nextHandler := range handlers {
				if nextHandler == nil {
					c.Next()
					continue
				} else {
					c.nextHandler = nextHandler
				}
			}

			for _, f := range h.Filters {
				switch t := f.(type) {
				case Filter:
					filtersResult = t.Check(c.Update)
				case func(u *objects.Update) bool:
					filtersResult = t(c.Update)
				default:
				}
			}

			if !filtersResult {
				c.Next()
				continue
			}
			c.GetCurrent()(c)
			c.Reset()
		}
	}
}

// Register could work in two modes
//
// 1. registers a filters, and handlers, this mode will register chain in the end of function
// 2. register handlerchain instantly
func (ho *HandlerObj) Register(callbacks ...interface{}) {
	var handler *HandlerType
	var partial bool

	handler = &HandlerType{}

	for _, elem := range callbacks {
		ho.mu.Lock()
		switch res := elem.(type) {
		case func(*objects.Update) error:
		case Filter:
			partial = true
			handler.AddFilters(res)

		case HandlerFunc:
			partial = true
			handler.Add(res)

		case *HandlerType:
		case HandlerType:
			ho.handlers = append(ho.handlers, &res)
		}
		ho.mu.Unlock()
	}
	if !partial {
		ho.handlers = append(ho.handlers, handler)
	}
}
