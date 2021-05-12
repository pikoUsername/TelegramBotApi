package dispatcher

import (
	"errors"

	"github.com/pikoUsername/tgp/objects"
)

type MiddlewareFunc func(update *objects.Update, handler HandlerType)

// Middleware is interface, default realization is DefaultMiddleware
type MiddlewareManager interface {
	Trigger(update *objects.Update, handler HandlerType)
	Register(middlewares ...MiddlewareFunc) // for many middleware add
	Unregister(middleware *MiddlewareFunc) (*MiddlewareFunc, error)
}

type DefaultMiddlewareManager struct {
	middlewares []*MiddlewareFunc
	dp          *Dispatcher
}

// NewDMiddlewareManager creates a DefaultMiddlewareManager, and return
func NewDMiddlewareManager(dp *Dispatcher) *DefaultMiddlewareManager {
	return &DefaultMiddlewareManager{
		dp: dp,
	}
}

// Trigger triggers special type of middlewares
// have three middleware types: pre, process, post
// We can register a middleware using Register Middleware
func (dmm *DefaultMiddlewareManager) Trigger(upd *objects.Update, handler HandlerType) {
	for _, cb := range dmm.middlewares {
		c := *cb
		c(upd, handler)
	}
}

// Register ...
func (dmm *DefaultMiddlewareManager) Register(md ...MiddlewareFunc) {
	// Transoforming Objects to Pointers
	// It s obuvious is bad code,
	// and maybe in golang libs exists func to make same, but more efficient
	var obj []*MiddlewareFunc

	for _, o := range md {
		obj = append(obj, &o)
	}

	dmm.middlewares = append(dmm.middlewares, obj...)
}

// Unregister a middleware
func (dmm *DefaultMiddlewareManager) Unregister(md *MiddlewareFunc) (*MiddlewareFunc, error) {
	// Checking for memory address
	for i, m := range dmm.middlewares {
		if m == md {
			// removing from list
			dmm.middlewares = append(dmm.middlewares[:i-1], dmm.middlewares[i:]...)
			return m, nil
		}
	}
	return nil, errors.New("this function not in middlewares")
}
