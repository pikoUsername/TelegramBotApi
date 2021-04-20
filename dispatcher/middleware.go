package dispatcher

// MiddlewareType is typeof callbacks
type MiddlewareType func(interface{}, *HandlerType)

// Middleware is interface, default realization is DefaultMiddleware
type Middleware interface {
	Trigger(obj interface{}, handler HandlerType) error
	Register(MiddlewareType) error
	Unregister(string) (*MiddlewareType, error)
	GetCallbacks() []MiddlewareType // for iteration
}

type MiddlewareDefault struct {
	Callbacks []MiddlewareType
}

func (md *MiddlewareDefault) Trigger(obj interface{}, handler HandlerType) error {
	for _, cb := range md.Callbacks {
		cb(obj, &handler)
	}
	return nil
}
