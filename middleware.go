package tgp

import (
	"errors"
	"reflect"

	"github.com/pikoUsername/tgp/objects"
)

// MiddlewareFunc is generic,
// first type, is pre-process middleware
// func(*objects.Update)
// second type is process middleware
// func(*objects.Update) bool/error
// and last type, third is post-process middleware
// func(objects.Update)
// passing just copy of update, bc pass by ptr havnot got any sense
type MiddlewareFunc interface{}

type PreMiddleware func(*objects.Update)
type PostMiddleware func(*objects.Update)
type ProcessMiddleware func(*Bot, *objects.Update) error

type triggerType map[string]func(c interface{}, upd *objects.Update, bot *Bot) error

// Middleware is interface, default realization is DefaultMiddleware
type MiddlewareManager interface {
	Trigger(bot *Bot, update *objects.Update, typ string) error
	Register(middlewares ...MiddlewareFunc) // for many middleware add
	Unregister(middleware *MiddlewareFunc) (*MiddlewareFunc, error)
}

const (
	PREMIDDLEWARE     = "pre"
	PROCESSMIDDLEWARE = "process"
	POSTMIDDLEWARE    = "post"
)

// errors
var (
	MiddlewareTypeInvalid = errors.New("typ parameter of variable not in ['post', 'pre', 'process']")
	MiddleawreNotFound    = errors.New("passed middleware not found")
)

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

// getConvertErr creates err in the fly with template:
// "failed convert this {value which failed to convert} to {tried to convert}"
func getConvertErr(it interface{}, ito interface{}) error {
	ts := reflect.TypeOf(it).String()
	tos := reflect.TypeOf(ito).String()
	return errors.New("failed convert this " + ts + " to " + tos)
}

// preTriggerProcess ...
func preTriggerProcess(c interface{}, upd *objects.Update, bot *Bot) error {
	if cb, ok := c.(PreMiddleware); ok {
		cb(upd)
	}
	return getConvertErr(c, (*PreMiddleware)(nil))
}

// processTrigger ...
func processTrigger(c interface{}, upd *objects.Update, bot *Bot) error {
	if process_middleware_cb, ok := c.(ProcessMiddleware); ok {
		err := process_middleware_cb(bot, upd)
		if err != nil {
			return err
		}
	}
	return getConvertErr(c, (ProcessMiddleware)(nil))
}

// postTrigger ...
func postTrigger(c interface{}, upd *objects.Update, bot *Bot) error {
	if post_middleware_cb, ok := c.(PostMiddleware); ok {
		post_middleware_cb(upd)
	}
	return getConvertErr(c, (*PostMiddleware)(nil))
}

// Trigger triggers special type of middlewares
// have three middleware types: pre, process, post
// We can register a middleware using Register Middleware
func (dmm *DefaultMiddlewareManager) Trigger(bot *Bot, upd *objects.Update, typ string) error {
	trigger_map := triggerType{
		PREMIDDLEWARE:     preTriggerProcess,
		PROCESSMIDDLEWARE: processTrigger,
		POSTMIDDLEWARE:    postTrigger,
	}

	for _, cb := range dmm.middlewares {
		c := *cb

		if val, ok := trigger_map[typ]; ok {
			err := val(c, upd, bot)
			if err != nil {
				return err
			}
		} else {
			return MiddlewareTypeInvalid
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
	// Checking for memory address, its really bad idea, but variant with map, too huge
	// variant with struct, too huge, and for middlewares store no need to use special structs
	for i, m := range dmm.middlewares {
		if m == md {
			// removing from list
			dmm.middlewares = append(dmm.middlewares[:i-1], dmm.middlewares[i:]...)
			return m, nil
		}
	}
	return nil, MiddleawreNotFound
}

func (dmm *DefaultMiddlewareManager) UnregisterByIndex(i uint) {
	dmm.middlewares = append(dmm.middlewares[i:], dmm.middlewares[:i+1]...)
}
