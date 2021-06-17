package dispatcher

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/pikoUsername/tgp/objects"
)

// MiddlewareFunc is generic,
// first type, is pre-process middleware
// func(*objects.Update)
// second type is process middleware
// func(*objects.Update) bool
// and last type, third is post-process middleware
// func(objects.Update)
// passing just copy of update, bc pass by ptr havnot got any sense
type MiddlewareFunc interface{}

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
	var err_mes string

	for _, cb := range dmm.middlewares {
		c := *cb
		pre_middlewre_cb, ok := c.(func(*objects.Update))

		if ok {
			pre_middlewre_cb(upd)

			continue
		} else {
			err_mes = "func(*objects.Update)"
		}
		process_middleware_cb, ok := c.(func(*objects.Update) error)

		if ok {
			err := process_middleware_cb(upd)
			if err != nil {
				return err
			}

			continue
		} else {
			process_middleware_cb, ok := c.(func(*objects.Update) bool)

			if ok {
				process_middleware_cb(upd)

				continue
			} else {
				err_mes = "func(*objects.Update) error / bool"
			}
		}
		post_middleware_cb, ok := c.(func(objects.Update))

		if ok {
			post_middleware_cb(*upd)

			continue
		} else {
			err_mes = "func(objects.Update)"
		}
		return errors.New("Failed convert this " + fmt.Sprintln(reflect.TypeOf(c)) + " to " + err_mes)
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
