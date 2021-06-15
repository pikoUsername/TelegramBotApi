package dispatcher

import (
	"errors"

	"github.com/pikoUsername/tgp/objects"
)

type MiddlewareFunc func(update *objects.Update) error

// Middleware is interface, default realization is DefaultMiddleware
type MiddlewareManager interface {
	Trigger(update *objects.Update) error
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
func (dmm *DefaultMiddlewareManager) Trigger(upd *objects.Update) error {
	for _, cb := range dmm.middlewares {
		c := *cb

		err := c(upd)
		if err != nil {
			return err
		}
	}

	return nil
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
