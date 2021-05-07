package dispatcher

import "github.com/pikoUsername/tgp/objects"

type MiddlewareFunc func(*objects.Update, HandlerFunc)

// Middleware is interface, default realization is DefaultMiddleware
type MiddlewareManager interface {
	Trigger(*objects.Update, HandlerFunc)
	Register(...MiddlewareFunc) // for many middleware add
	Unregister(MiddlewareFunc) (*MiddlewareFunc, error)
}

type DefaultMiddlewareManager struct {
	middlewares []MiddlewareFunc
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
func (dmm *DefaultMiddlewareManager) Trigger(upd *objects.Update, handler HandlerFunc) {
	for _, cb := range dmm.middlewares {
		cb(upd, handler)
	}
}

// Register ...
func (dmm *DefaultMiddlewareManager) Register(md ...MiddlewareFunc) {
	dmm.middlewares = append(dmm.middlewares, md...)
}

// Unregister a middleware
func (dmm *DefaultMiddlewareManager) Unregister(md MiddlewareFunc) (*MiddlewareFunc, error) {
	return &md, nil // magic!
}
