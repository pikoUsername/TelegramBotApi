package dispatcher

import "github.com/pikoUsername/tgp/objects"

type MiddlewareFunc func(*objects.Update, HandlerType)

// Middleware is interface, default realization is DefaultMiddleware
type MiddlewareManager interface {
	Trigger(update *objects.Update, handler HandlerType)
	Register(middlewares ...MiddlewareFunc) // for many middleware add
	Unregister(middleware MiddlewareFunc) (*MiddlewareFunc, error)
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
func (dmm *DefaultMiddlewareManager) Trigger(upd *objects.Update, handler HandlerType) {
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
