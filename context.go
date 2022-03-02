package tgp

import (
	"encoding/json"
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

// Context object used in middlewares, handlers
// middlewares can write some data for interact with handler
// p.s idea taken from gin sources
type Context struct {
	*objects.Update

	Bot      *Bot
	Storage  storage.Storage
	Markdown Markdown

	data     map[string]interface{}
	index    int
	cursor   int
	handlers []*HandlerType

	calledErrors []error
	mu           sync.Mutex

	hasDone chan struct{}
}

// Context.Set just set ctxVar to key in data context
func (ctx *Context) Set(key string, contextVar interface{}) {
	ctx.mu.Lock()
	ctx.data[key] = contextVar
	ctx.mu.Unlock()
}

// Context.Get do not notify about error, error will be ignored
func (ctx *Context) Get(key string) (v interface{}, ok bool) {
	v, ok = ctx.data[key]
	return
}

// MustGet Same as Get, but dont checks a existing,
// instead call Fatal method
func (ctx *Context) MustGet(key string) (v interface{}) {
	if v, ok := ctx.Get(key); v != nil && ok {
		return v
	}
	ctx.Fatal(tgpErr.New("Key " + key + " does not exists"))
	return
}

// Next calls next handler, and increment cursor
func (ctx *Context) Next() {
	if ctx.index >= AcceptIndex {
		if ctx.cursor >= len(ctx.handlers) {
			return
		}
		hand := ctx.GetCurrent()

		ctx.mu.Lock()
		ctx.cursor++
		ctx.mu.Unlock()

		ctx.call(hand)
	}
}

func (ctx *Context) call(hand *HandlerType) {
	if len(hand.filters) == 0 || checkFilters(hand.filters, ctx.Update) {
		hand.GetHandler()(ctx)
		ctx.hasDone <- struct{}{}
	}
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.hasDone
}

// Returns Context cursor
func (ctx *Context) Cursor() int {
	return ctx.cursor
}

// Err returns error which raised in pervious handlers
func (ctx *Context) GetErrors() []error {
	return ctx.calledErrors
}

func (ctx *Context) GetCurrent() *HandlerType {
	return ctx.handlers[ctx.cursor]
}

// Abort sets index to AbortIndex
func (ctx *Context) Abort() {
	ctx.index = AbortIndex
}

// AbortWithError ...
func (ctx *Context) AbortWithError(err error) []error {
	ctx.Abort()
	ctx.mu.Lock()
	ctx.calledErrors = append(ctx.calledErrors, err)
	ctx.mu.Unlock()
	return ctx.calledErrors
}

// Error adds to errors list errors from arguments
func (ctx *Context) Error(s ...interface{}) error {
	err := errors.New(fmt.Sprintln(s...))
	ctx.mu.Lock()
	ctx.calledErrors = append(ctx.calledErrors, err)
	ctx.mu.Unlock()
	return err
}

func (ctx *Context) Errorf(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	ctx.mu.Lock()
	ctx.calledErrors = append(ctx.calledErrors, err)
	ctx.mu.Unlock()
	return err
}

func (ctx *Context) Fatalf(format string, args ...interface{}) error {
	ctx.Abort()
	return ctx.Errorf(format, args...)
}

// Fatal calls Abort method, and do same thing as Error
func (ctx *Context) Fatal(s ...interface{}) error {
	ctx.Abort()
	return ctx.Error(s...)
}

// InputFile ...
func (ctx *Context) InputFile(name, path string) (*objects.InputFile, error) {
	return objects.NewInputFile(path, name)
}

// IsMessageToMe returns true if message directed to this bot.
func (ctx *Context) IsMessageToMe(message *objects.Message) bool {
	return strings.Contains(message.Text, "@"+ctx.Bot.Me.Username)
}

// Sends message request, which must return Message object.
// if request type is not correct, will return error
func (ctx *Context) Send(config Configurable) (*objects.Message, error) {
	return ctx.Bot.Send(config)
}

// Reply to this context object
func (ctx *Context) Reply(config Configurable) (*objects.Message, error) {
	var upd = ctx.Update
	var chat *objects.Chat

	if upd.EditedMessage != nil {
		chat = upd.EditedMessage.Chat
	} else if upd.ChannelPost != nil {
		chat = upd.ChannelPost.Chat
	} else if upd.Message != nil {
		chat = upd.Message.Chat
	} else {
		return &objects.Message{}, tgpErr.New("Update is empty")
	}
	chat_id_str := strconv.FormatInt(chat.ID, 10)

	// code duplication
	switch conf := config.(type) {
	case FileableConf:
		params, err := conf.params()
		params["chat_id"] = chat_id_str
		if err != nil {
			return &objects.Message{}, err
		}

		method := config.method()
		resp, err := ctx.Bot.UploadFile(method, params, conf.getFiles()...)
		if err != nil {
			return &objects.Message{}, err
		}

		var message *objects.Message
		json.Unmarshal(resp.Result, &message)
		return message, nil
	case Configurable:
		v, err := conf.values()
		v.Set("chat_id", chat_id_str)
		if err != nil {
			return &objects.Message{}, err
		}
		if v.Get("parse_mode") == "" {
			v.Set("parse_mode", ctx.Bot.ParseMode)
		}
		resp, err := ctx.Bot.Request(conf.method(), v)

		if err != nil {
			return &objects.Message{}, err
		}
		var msg objects.Message
		json.Unmarshal(resp.Result, &msg)
		return &msg, nil
	}
	return &objects.Message{}, tgpErr.New("config is not correct")
}

// SetState set a state which passed for a current user in current chat
// works only in handler, or in middleware, nor outside
func (ctx *Context) SetState(state *fsm.State) error {
	cid, uid := extractIds(ctx.Update)
	return ctx.Storage.SetState(cid, uid, state.GetFullState())
}

func (ctx *Context) GetState() (*fsm.State, error) {
	var state, group string
	cid, uid := extractIds(ctx.Update)
	st, err := ctx.Storage.GetState(cid, uid)
	if err != nil {
		return &fsm.State{}, err
	}
	s := strings.Split(st, ":")
	if len(s) < 2 {
		return &fsm.State{}, tgpErr.New("Uncorrect state string format")
	}
	state, group = s[0], s[1]
	return &fsm.State{State: state, GroupState: group}, nil
}

// ResetState reset state for current user, and current chat
func (ctx *Context) ResetState() error {
	cid, uid := extractIds(ctx.Update)
	return ctx.Storage.SetState(cid, uid, fsm.DefaultState.GetFullState())
}

func (ctx *Context) Reset() {
	ctx.calledErrors = ctx.calledErrors[:0]
	ctx.data = nil
	ctx.handlers = ctx.handlers[:0]
	ctx.Update = nil
}
