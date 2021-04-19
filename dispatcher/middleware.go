package dispatcher

import "github.com/pikoUsername/tgp/bot"

// MiddlewareType is typeof callbacks
type MiddlewareType func(interface{}, *func(interface{}, bot.Bot))

// Middleware is interface, default realization is DefaultMiddleware
type Middleware interface {
	Trigger(string) error
	Register(string, MiddlewareType) error
	Unregister(string) (*MiddlewareType, error)
	GetCallbacks() []MiddlewareType // for iteration
}
