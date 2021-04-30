package dispatcher

import "github.com/pikoUsername/tgp/objects"

// MiddlewareType is typeof callbacks
// Middleware is one type, but you can make various middlewares,
// and activate command in any place of your middleware, pshe
type MiddlewareFunc func(*objects.Update, HandlerFunc)

// Middleware is interface, default realization is DefaultMiddleware
type MiddlewareManager interface {
	Trigger(*objects.Update, HandlerFunc)
	Register(MiddlewareType)
	Unregister(MiddlewareType) (*MiddlewareType, error)
}

type DefaultMiddlewareManager struct {
	middlewares []MiddlewareType
	dp          *Dispatcher
}

// NewDMiddlewareManager creates a DefaultMiddlewareManager, and return
func NewDMiddlewareManager(dp *Dispatcher) *DefaultMiddlewareManager {
	return &DefaultMiddlewareManager{
		dp: dp,
	}
}

// Trigger uses for trigger all middlewares
// I write this againm you must to call handler in the middleware function
func (dmm *DefaultMiddlewareManager) Trigger(upd *objects.Update, handler HandlerFunc) {
	for _, cb := range dmm.middlewares {
		cb(upd, handler)
	}
}

// Register ...
func (dmm *DefaultMiddlewareManager) Register(md MiddlewareType) {
	dmm.middlewares = append(dmm.middlewares, md)
}

// Unregister a middleware
func (dmm *DefaultMiddlewareManager) Unregister(md MiddlewareType) (*MiddlewareType, error) {
	return &md, nil // magic!
}
