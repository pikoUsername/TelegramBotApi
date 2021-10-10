package tgp

import (
	"errors"
	"fmt"
	"sync"

	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

const (
	AbortIndex = iota
	ContinueIndex
)

// dataContext have no mechanisms for syncronize data
type dataContext map[string]interface{}

func (dc dataContext) Get(key string) (v interface{}, ok bool) {
	v, ok = dc[key]
	return
}

func (dc dataContext) Set(key string, val interface{}) {
	dc[key] = val
}

// Context object used in middlewares, handlers
// middlewares can write some data for interact with handler
// p.s idea taken from gin sources
type Context struct {
	*objects.Update

	Bot      *Bot
	Storage  storage.Storage
	data     dataContext
	index    int32
	cursor   int32
	handlers []HandlerFunc

	// nextHandler maybe a just handler, or middleware
	nextHandler  func(ctx *Context)
	calledErrors []error

	mu sync.Mutex
}

// Context.Set just set ctxVar to key in data context
func (c *Context) Set(key string, contextVar interface{}) {
	c.mu.Lock()
	c.data.Set(key, contextVar)
	c.mu.Unlock()
}

// Context.Get do not notify about error, error will be ignored
func (c *Context) Get(key string) (v interface{}, ok bool) {
	v, ok = c.data.Get(key)
	return
}

// MustGet Same as Get, but dont checks a existing,
// instead call Fatal method
func (c *Context) MustGet(key string) (v interface{}) {
	if v, ok := c.Get(key); v != nil && ok {
		return v
	}
	c.Fatal(tgpErr.New("Key " + key + " does not exists"))
	return
}

// Next calls next handler
func (c *Context) Next() {
	if c.index >= ContinueIndex {
		c.nextHandler(c)
		c.cursor++
	}
}

func (c *Context) GetCurrent() HandlerFunc {
	return c.handlers[c.cursor]
}

// Abort sets index to AbortIndex
func (c *Context) Abort() {
	c.index = AbortIndex
}

// AbortWithError ...
func (c *Context) AbortWithError(err error) []error {
	c.Abort()
	c.mu.Lock()
	c.calledErrors = append(c.calledErrors, err)
	c.mu.Unlock()
	return c.calledErrors
}

// Error adds to errors list errors from arguments
func (c *Context) Error(s ...interface{}) []error {
	err := errors.New(fmt.Sprintln(s...))
	c.mu.Lock()
	c.calledErrors = append(c.calledErrors, err)
	c.mu.Unlock()
	return c.calledErrors
}

// Fatal calls Abort method, and do same thing as Error
func (c *Context) Fatal(s ...interface{}) []error {
	c.Abort()
	return c.Error(s...)
}

func (c *Context) Reset() {
	c.index = -1
	// c.data = nil
	c.calledErrors = c.calledErrors[:0] // what?
	c.nextHandler = nil
}

func (c *Context) InputFile(name, path string) (*objects.InputFile, error) {
	var file *objects.InputFile

	if name == "" || path == "" {
		return file, tgpErr.New("Name and/or path arguments is unfilled ")
	}
	return objects.NewInputFile(path, name)
}

// func (c *Context) Deadline() (deadline time.Time, ok bool) {}
// func (c *Context) Done() <-chan struct{}                   {}
// func (c *Context) Err() error                              {}
// func (c *Context) Value(key interface{}) interface{}       {}
