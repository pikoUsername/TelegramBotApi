package dispatcher

// MiddlewareType is typeof callbacks
// Middleware is one type, but you can make various middlewares,
// and activate command in any place of your middleware, pshe
type MiddlewareType func(*Context, *HandlerType) interface{}

// Middleware is interface, default realization is DefaultMiddleware
type MiddlewareManager interface {
	Trigger(*Context, *HandlerType)
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
func (dmm *DefaultMiddlewareManager) Trigger(ctx *Context, handler *HandlerType) {
	for _, cb := range dmm.middlewares {
		cb(ctx, handler)
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
