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
//
// Not concurrency safe!
type HandlerType struct {
	Callbacks []HandlerFunc
	Filters   []interface{}
}

func (h *HandlerType) Add(handlers ...HandlerFunc) {
	if len(handlers) != 0 {
		h.Callbacks = append(h.Callbacks, handlers...)
	}
}

func (h *HandlerType) AddFilters(filters ...interface{}) {
	if len(filters) != 0 {
		h.Filters = append(h.Filters, filters...)
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
			if len(c.calledErrors) != 0 {

			}
			ho.handlerPos = *(*uint16)(unsafe.Pointer(&i))
			nextIndex = i + 1
			if nextIndex > len(ho.handlers) {
				nextIndex = i
			}
			handlers := ho.handlers[nextIndex].Callbacks
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
					filtersResult = false
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

// Register, append to Callbacks, e.g handler functions
func (ho *HandlerObj) Register(f HandlerFunc, filters ...interface{}) {
	ht := HandlerType{
		Callbacks: []HandlerFunc{f},
		Filters:   filters,
	}

	ho.mu.Lock()
	ho.handlers = append(ho.handlers, &ht)
	ho.mu.Unlock()
}

func (ho *HandlerObj) RegisterChain(f *HandlerType) {
	ho.mu.Lock()
	ho.handlers = append(ho.handlers, f)
	ho.mu.Unlock()
}
