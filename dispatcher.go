package tgp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pikoUsername/tgp/fsm"
	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

// Dispatcher's purpose is run bot, and comfortable pipeline
// Bot struct uses as API wrapper
// Dispatcher uses as Bot starter
// Another level of abstraction
type Dispatcher struct {
	Bot *Bot
	// Handlers
	MessageHandler         *HandlerObj
	CallbackQueryHandler   *HandlerObj
	ChannelPostHandler     *HandlerObj
	PollHandler            *HandlerObj
	ChatMemberHandler      *HandlerObj
	PollAnswerHandler      *HandlerObj
	MyChatMemberHandler    *HandlerObj
	ChatJoinRequestHandler *HandlerObj

	// Storage interface
	Storage storage.Storage
	logger  StdLogger

	// for FSM usage
	currentUpdate *objects.Update

	// If you want to add onshutdown function
	// just append to this object, :P
	OnWebhookShutdown []OnStartAndShutdownFunc
	OnPollingShutdown []OnStartAndShutdownFunc
	OnWebhookStartup  []OnStartAndShutdownFunc
	OnPollingStartup  []OnStartAndShutdownFunc

	Welcome bool
	polling bool
	webhook bool

	closeChan   chan struct{}
	functionsWG *sync.WaitGroup
}

var (
	ErrorTypeAssertion = tgpErr.New("impossible to do type assertion to this callback")
	ErrorConflictModes = tgpErr.New("enabled two conflicting modes at the same time, polling and webhook")
)

type OnStartAndShutdownFunc func(dp *Dispatcher)

// NewDispathcer get a new Dispatcher with default values
func NewDispatcher(bot *Bot, storage storage.Storage) *Dispatcher {
	dp := &Dispatcher{
		Bot:         bot,
		Storage:     storage,
		Welcome:     true,
		logger:      bot.logger, // pointer to logger
		functionsWG: &sync.WaitGroup{},
		closeChan:   make(chan struct{}),
	}

	dp.MessageHandler = NewHandlerObj()
	dp.CallbackQueryHandler = NewHandlerObj()
	dp.ChannelPostHandler = NewHandlerObj()
	dp.ChatMemberHandler = NewHandlerObj()
	dp.PollHandler = NewHandlerObj()
	dp.PollAnswerHandler = NewHandlerObj()
	dp.ChannelPostHandler = NewHandlerObj()

	return dp
}

// OnConfig using as argument for OnStartup, OnShutdown methods
// You can add multiple functions to startup, or shutdown mthds
// Example:
// c := &OnConfig{}
// // could be added a multiple functions in one call
// c.Add(func(...) {})
// dp.OnStartup(c)
type OnConfig struct {
	Polling bool
	Webhook bool
	cb      []OnStartAndShutdownFunc
}

func (oc *OnConfig) Add(cb OnStartAndShutdownFunc) {
	oc.cb = append(oc.cb, cb)
}

func NewOnConf(cb OnStartAndShutdownFunc) *OnConfig {
	return &OnConfig{
		cb:      []OnStartAndShutdownFunc{cb},
		Webhook: true,
		Polling: true,
	}
}

func callListFuncs(funcs []OnStartAndShutdownFunc, dp *Dispatcher) {
	for _, cb := range funcs {
		dp.functionsWG.Add(1)
		go cb(dp)
	}
	dp.functionsWG.Wait()
}

// Config for start polling method
type StartPollingConfig struct {
	*GetUpdatesConfig
	Context      context.Context
	SkipUpdates  bool
	SafeExit     bool
	ResetWebhook bool
	ErrorSleep   uint
	Relax        time.Duration
	Timeout      time.Duration
}

func NewStartPollingConf(skip_updates bool) *StartPollingConfig {
	return &StartPollingConfig{
		GetUpdatesConfig: &GetUpdatesConfig{
			Timeout: 5,
		},
		Relax:        1 * time.Second, // 0.1
		ResetWebhook: false,
		ErrorSleep:   1,
		SkipUpdates:  skip_updates,
		SafeExit:     true,
		Timeout:      5 * time.Second,
	}
}

// StartWebhookConfig uses for start bot using webhooks
type StartWebhookConfig struct {
	*SetWebhookConfig
	Handler            http.Handler
	KeyFile            interface{}
	BotURL             string
	Address            string
	DropPendingUpdates bool
	SafeExit           bool
}

func NewStartWebhookConf(url string, address string) *StartWebhookConfig {
	return &StartWebhookConfig{
		BotURL:  url,
		Address: address,
	}
}

// ResetWebhook uses for reset webhook for telegram
func (dp *Dispatcher) ResetWebhook(check bool) error {
	if check {
		wi, err := dp.Bot.GetWebhookInfo()
		if err != nil {
			return err
		}
		if wi.URL == "" {
			return errors.New("url is nothing")
		}
	}
	_, err := dp.Bot.DeleteWebhook(&DeleteWebhookConfig{})
	return err
}

// ProcessOneUpdate processes only one comming update
func (dp *Dispatcher) ProcessOneUpdate(upd *objects.Update) error {
	local_ctx := dp.Context(upd)

	if upd.Message != nil {
		dp.MessageHandler.Trigger(local_ctx)
	} else if upd.CallbackQuery != nil {
		dp.CallbackQueryHandler.Trigger(local_ctx)
	} else if upd.ChannelPost != nil {
		dp.ChannelPostHandler.Trigger(local_ctx)
	} else if upd.Poll != nil {
		dp.PollHandler.Trigger(local_ctx)
	} else if upd.PollAnswer != nil {
		dp.PollAnswerHandler.Trigger(local_ctx)
	} else if upd.ChatMember != nil {
		dp.ChatMemberHandler.Trigger(local_ctx)
	} else if upd.MyChatMember != nil {
		dp.MyChatMemberHandler.Trigger(local_ctx)
	} else if upd.ChatJoinRequest != nil {
		dp.ChatJoinRequestHandler.Trigger(local_ctx)
	} else {
		text := "detected not supported type of updates, seems like telegram bot api updated before this package updated"
		return tgpErr.New(text)
	}

	return nil
}

// SkipUpdates skip comming updates, sending to telegram servers
func (dp *Dispatcher) SkipUpdates() (err error) {
	_, err = dp.Bot.GetUpdates(&GetUpdatesConfig{
		Offset:  -1,
		Timeout: 1,
	})
	return
}

func (dp *Dispatcher) Context(upd *objects.Update) *Context {
	return &Context{
		Update:  upd,
		data:    dataContext{},
		index:   AcceptIndex,
		Bot:     dp.Bot,
		Storage: dp.Storage,
		mu:      sync.Mutex{},
	}
}

// SetState set a state which passed for a current user in current chat
// works only in handler, or in middleware, nor outside
func (dp *Dispatcher) SetState(state *fsm.State) error {
	if dp.currentUpdate != nil {
		cid, uid := getUidAndCidFromUpd(dp.currentUpdate)
		return dp.Storage.SetState(cid, uid, state.GetFullState())
	}
	return nil
}

// ResetState reset state for current user, and current chat
func (dp *Dispatcher) ResetState() error {
	if dp.currentUpdate != nil {
		cid, uid := getUidAndCidFromUpd(dp.currentUpdate)
		return dp.Storage.SetState(cid, uid, fsm.DefaultState.GetFullState())
	}
	return nil
}

// ========================================
//   Startup and Shutdown related methods
// ========================================

// Shutdown calls when you enter ^C(which means SIGINT)
// And SafeExit catch it, before you exit
func (dp *Dispatcher) shutdownPolling() {
	if len(dp.OnPollingShutdown) == 0 {
		return
	}
	callListFuncs(dp.OnPollingShutdown, dp)
}

// startUpPolling function, iterate over a callbacks from OnStartupCallbacks
// Calls in StartPolling function
func (dp *Dispatcher) startupPolling() {
	go dp.welcome()
	if len(dp.OnPollingStartup) == 0 {
		return
	}
	callListFuncs(dp.OnPollingStartup, dp)
}

// shutdownWebhook method, iterate over a callbacks from OnWebhookShutdown
func (dp *Dispatcher) shutdownWebhook() {
	if len(dp.OnWebhookShutdown) == 0 {
		return
	}
	callListFuncs(dp.OnWebhookShutdown, dp)
}

// startupPolling method, iterate over a callbacks from OnWebhookStartup
func (dp *Dispatcher) startupWebhook() {
	go dp.welcome()
	if len(dp.OnWebhookStartup) == 0 {
		return
	}
	callListFuncs(dp.OnWebhookStartup, dp)
}

// Onstartup method append to OnStartupCallbaks a callbacks
// Using pointers bc cant unregister function using copy of object
// And golang doesnot support generics, and type equals
func (dp *Dispatcher) OnStartup(c *OnConfig) {
	if !c.Webhook && !c.Polling {
		dp.logger.Println("this expression have not got any effect")
	}

	if c.Webhook {
		dp.OnWebhookStartup = append(dp.OnWebhookStartup, c.cb...)
	}
	if c.Polling {
		dp.OnPollingStartup = append(dp.OnPollingStartup, c.cb...)
	}
}

// OnShutdown method using for register OnShutdown callbacks
// Same code like OnStartup
func (dp *Dispatcher) OnShutdown(c *OnConfig) {
	if !c.Webhook && !c.Polling {
		dp.logger.Println("!polling and !webhook expression have not got any effect")
	}

	if c.Webhook {
		dp.OnWebhookShutdown = append(dp.OnWebhookShutdown, c.cb...)
	}
	if c.Polling {
		dp.OnPollingShutdown = append(dp.OnPollingShutdown, c.cb...)
	}
}

func (dp *Dispatcher) Start() {
	if dp.polling {
		dp.startupPolling()
	}
	if dp.webhook {
		dp.startupWebhook()
	}
}

func (dp *Dispatcher) runShutDown() {
	if dp.polling {
		dp.shutdownPolling()
	}
	if dp.webhook {
		dp.shutdownWebhook()
	}
}

// SafeExit method uses for notify about exit from program
// Thanks: https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
func (dp *Dispatcher) SafeExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		dp.ShutDown()
		os.Exit(0)
	}()
}

// ShutDownDP calls ResetWebhook
func (dp *Dispatcher) ShutDown() {
	dp.logger.Println("Stop polling!")
	dp.ResetWebhook(true)
	dp.Storage.Close()
	dp.runShutDown()
}

func (dp *Dispatcher) welcome() {
	if dp.Welcome {
		_, err := dp.Bot.GetMe()
		if err != nil {
			dp.logger.Printf("tgp: %s", err)
			return
		}
		dp.logger.Println("Bot: ", dp.Bot.Me.Username)
	}
}

// =========================================
//    Polling and webhook related methods
// =========================================

// GetUpdatesChan makes getUpdates request to telegram servers
// sends update to updates channel
// Time.Sleep here for stop goroutine for a c.Relax time
//
// yeah it bad, and works only on crutches, but works
func (dp *Dispatcher) MakeUpdatesChan(c *StartPollingConfig, ch chan *objects.Update) {
	go func() {
		for {
			if c.Relax != 0 {
				time.Sleep(c.Relax)
			}

			updates, err := dp.Bot.GetUpdates(c.GetUpdatesConfig)
			if err != nil {
				dp.logger.Println(err.Error())
				dp.logger.Println("Error with getting updates")
				time.Sleep(time.Duration(c.ErrorSleep))

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= c.Offset {
					c.Offset = update.UpdateID + 1
					ch <- update
				}
			}
		}
	}()
}

// ProcessUpdates iterates <-chan *objects.Update
func (dp *Dispatcher) ProcessUpdates(ch <-chan *objects.Update) error {
	cherr := make(chan error, 1)

	for {
		upd := <-ch
		if upd == nil {
			continue
		}
		go func() { cherr <- dp.ProcessOneUpdate(upd) }()

		select {
		case err := <-cherr:
			if err != nil {
				return err
			}
		case <-dp.closeChan:
			dp.ShutDown()
			return nil
		}
	}
}

// StartPolling check out to comming updates
// If yes, Telegram Get to your bot a Update
// Using GetUpdates method in Bot structure
// GetUpdates config using for getUpdates method
func (dp *Dispatcher) StartPolling(c *StartPollingConfig) error {
	var err error
	if dp.webhook {
		return ErrorConflictModes
	}

	dp.polling = true
	dp.Start()
	if c.SafeExit {
		dp.SafeExit()
	}
	if c.ResetWebhook {
		err = dp.ResetWebhook(true)
		if err != nil {
			return err
		}
	}

	if c.SkipUpdates {
		err = dp.SkipUpdates()
		if err != nil {
			return err
		}
	}

	ch := make(chan *objects.Update)
	dp.MakeUpdatesChan(c, ch)

	dp.ProcessUpdates(ch)
	return nil
}

// MakeWebhookChan adds a http Handler with c.BotURL path
func (dp *Dispatcher) MakeWebhookChan(c *StartWebhookConfig, ch chan *objects.Update) {
	http.HandleFunc(c.BotURL, func(wr http.ResponseWriter, req *http.Request) {
		update, err := requestToUpdate(req)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			wr.WriteHeader(http.StatusBadRequest)
			wr.Header().Set("Content-Type", "application/json")
			wr.Write(errMsg)
			return
		}

		ch <- update
	})
}

// StartWebhook method registers BotUrl uri a function which handles every comming update
// Using In Pair of SetWebhook method
// Startup method executes after SetWebhook method call
//
// NOTE: you should to add a webhook close callback function, using OnShutdown
func (dp *Dispatcher) StartWebhook(c *StartWebhookConfig) error {
	if dp.polling {
		panic(ErrorConflictModes)
	}
	_, err := dp.Bot.SetWebhook(c.SetWebhookConfig)
	if err != nil {
		return err
	}
	dp.webhook = true
	dp.Start()
	if c.SafeExit {
		dp.SafeExit()
	}
	http.HandleFunc(c.BotURL, func(wr http.ResponseWriter, req *http.Request) {
		update, err := requestToUpdate(req)
		if err != nil {
			WriteRequestError(wr, err)
			return
		}

		dp.currentUpdate = update
		err = dp.ProcessOneUpdate(update) // will run in stock mode
		if err != nil {
			WriteRequestError(wr, err)
			return
		}
	})
	certPath, err := guessFileName(c.Certificate)
	if err != nil {
		return err
	}
	keyfile, err := guessFileName(c.KeyFile)
	if err != nil {
		return err
	}
	return http.ListenAndServeTLS(c.Address, certPath, keyfile, c.Handler)
}
