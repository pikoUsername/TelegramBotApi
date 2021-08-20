package tgp

import (
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

const (
	PREMIDDLEWARE     = "pre"
	PROCESSMIDDLEWARE = "process"
	POSTMIDDLEWARE    = "post"
)

// errors
var (
	MiddlewareTypeInvalid = Errors.New("typ parameter of variable not in ['post', 'pre', 'process']")
	MiddlewareNotFound    = Errors.New("passed middleware not found")
	MiddlewareIncorrect   = Errors.New("passed function is not function type")
)

type DefaultMiddlewareManager struct {
	middlewares []MiddlewareFunc
}

// NewDMiddlewareManager creates a DefaultMiddlewareManager, and return
func NewMiddlewareManager(dp *Dispatcher) *DefaultMiddlewareManager {
	return &DefaultMiddlewareManager{}
}

// convertErr creates err in the fly with template:
// "failed convert this {value which failed to convert} to {tried to convert}"
func convertErr(it interface{}, ito interface{}) error {
	ts := reflect.TypeOf(it).String()
	tos := reflect.TypeOf(ito).String()
	return Errors.New("failed convert this " + ts + " to " + tos)
}

// preTriggerProcess ...
func preTriggerProcess(c interface{}, upd *objects.Update, bot *Bot) error {
	if cb, ok := c.(PreMiddleware); ok {
		cb(upd)
	}
	return convertErr(c, (*PreMiddleware)(nil))
}

// processTrigger ...
func processTrigger(c interface{}, upd *objects.Update, bot *Bot) error {
	if process_middleware_cb, ok := c.(ProcessMiddleware); ok {
		err := process_middleware_cb(bot, upd)
		if err != nil {
			return err
		}
	}
	return convertErr(c, (ProcessMiddleware)(nil))
}

// postTrigger ...
func postTrigger(c interface{}, upd *objects.Update, bot *Bot) error {
	if post_middleware_cb, ok := c.(PostMiddleware); ok {
		post_middleware_cb(upd)
	}
	return convertErr(c, (*PostMiddleware)(nil))
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
		if val, ok := trigger_map[typ]; ok {
			err := val(cb, upd, bot)
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
	dmm.middlewares = append(dmm.middlewares, md...)
}

// Unregister a middleware from middleware list
// thx: https://stackoverflow.com/questions/34901307/how-to-compare-2-functions-in-go/34901677
func (dmm *DefaultMiddlewareManager) Unregister(md MiddlewareFunc) (MiddlewareFunc, error) {
	t := reflect.TypeOf(md)
	if t.Kind() != reflect.Func {
		return nil, MiddlewareIncorrect
	}
	var s, s2 uintptr
	s = reflect.ValueOf(md).Pointer()

	for i, m := range dmm.middlewares {
		s2 = reflect.ValueOf(m).Pointer()
		if s == s2 {
			// removing from list
			i2 := i - 1
			if i2 < 0 {
				i2 = 0
			}
			dmm.middlewares = append(dmm.middlewares[:i2], dmm.middlewares[i2:]...)
			return m, nil
		}
	}
	return nil, MiddlewareNotFound
}

func (dmm *DefaultMiddlewareManager) UnregisterByIndex(i uint) {
	dmm.middlewares = append(dmm.middlewares[i:], dmm.middlewares[:i+1]...)
}
