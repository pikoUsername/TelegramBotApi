package dispatcher

import (
	"errors"

	"github.com/pikoUsername/tgp/objects"
)

type MiddlewareFunc func(update *objects.Update) bool

// Middleware is interface, default realization is DefaultMiddleware
type MiddlewareManager interface {
	Trigger(update *objects.Update) bool
	Register(middlewares ...MiddlewareFunc) // for many middleware add
	Unregister(middleware *MiddlewareFunc) (*MiddlewareFunc, error)
}

type DefaultMiddlewareManager struct {
	middlewares []*MiddlewareFunc
	dp          *Dispatcher
}

// NewDMiddlewareManager creates a DefaultMiddlewareManager, and return
func NewMiddlewareManager(dp *Dispatcher) *DefaultMiddlewareManager {
	dmm := &DefaultMiddlewareManager{
		dp: dp,
	}

	return dmm
}

// Trigger triggers special type of middlewares
// have three middleware types: pre, process, post
// We can register a middleware using Register Middleware
func (dmm *DefaultMiddlewareManager) Trigger(upd *objects.Update) bool {
	for _, cb := range dmm.middlewares {
		c := *cb

		b := c(upd)
		if !b {
			return false
		}
	}

	return true
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

func (dmm *DefaultMiddlewareManager) UnregisterByIndex(i uint) {
	dmm.middlewares = append(dmm.middlewares[i:], dmm.middlewares[:i+1]...)
}
