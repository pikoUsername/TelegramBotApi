package tgp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/pikoUsername/tgp/fsm"
	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

const (
	AbortIndex = iota
	AcceptIndex
)

// dataContext using for interact with each handlers
// Not thread safe
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
	index    int
	cursor   int
	handlers []*HandlerType

	calledErrors []error
	mu           sync.Mutex
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

// Err returns error which raised in pervious handlers
func (c *Context) Errors() []error {
	return c.calledErrors
}

// Next calls next handler, and increment cursor
func (c *Context) Next() {
	if c.index >= AcceptIndex {
		c.cursor++
		if c.cursor >= len(c.handlers)-1 {
			return
		}
		c.GetCurrent().apply(c)
	}
}

// Returns Context cursor
func (c *Context) Cursor() int {
	return c.cursor
}

func (c *Context) GetCurrent() *HandlerType {
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
func (c *Context) Error(s ...interface{}) error {
	err := errors.New(fmt.Sprintln(s...))
	c.mu.Lock()
	c.calledErrors = append(c.calledErrors, err)
	c.mu.Unlock()
	return err
}

func (c *Context) Errorf(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	c.mu.Lock()
	c.calledErrors = append(c.calledErrors, err)
	c.mu.Unlock()
	return err
}

func (c *Context) Fatalf(format string, args ...interface{}) error {
	c.Abort()
	return c.Errorf(format, args...)
}

// Fatal calls Abort method, and do same thing as Error
func (c *Context) Fatal(s ...interface{}) error {
	c.Abort()
	return c.Error(s...)
}

func (c *Context) InputFile(name, path string) (*objects.InputFile, error) {
	var file *objects.InputFile

	if name == "" || path == "" {
		return file, tgpErr.New("Name and/or path arguments is unfilled ")
	}
	return objects.NewInputFile(path, name)
}

// IsMessageToMe returns true if message directed to this bot.
func (ctx *Context) IsMessageToMe(message *objects.Message) bool {
	return strings.Contains(message.Text, "@"+ctx.Bot.Me.Username)
}

// Sends message request, which must return Message object.
// if request type is not correct, will return error
func (c *Context) Send(config Configurable) (*objects.Message, error) {
	return c.Bot.Send(config)
}

// Reply to this context object
func (c *Context) Reply(config Configurable) (*objects.Message, error) {
	v, _ := config.values()
	v.Set("chat_id", strconv.FormatInt(c.Message.Chat.ID, 10))
	return c.Send(config)
}

// SetState set a state which passed for a current user in current chat
// works only in handler, or in middleware, nor outside
func (ctx *Context) SetState(state *fsm.State) error {
	cid, uid := getUidAndCidFromUpd(ctx.Update)
	return ctx.Storage.SetState(cid, uid, state.GetFullState())
}

// ResetState reset state for current user, and current chat
func (ctx *Context) ResetState() error {
	cid, uid := getUidAndCidFromUpd(ctx.Update)
	return ctx.Storage.SetState(cid, uid, fsm.DefaultState.GetFullState())
}

// TODO:
// func (c *Context) Deadline() (deadline time.Time, ok bool) {}
// func (c *Context) Done() <-chan struct{}                   {}
// func (c *Context) Err() error                              {}
// func (c *Context) Value(key interface{}) interface{}       {}
