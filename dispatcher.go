package tgp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"time"

	"github.com/pikoUsername/tgp/fsm"
	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
	"github.com/pikoUsername/tgp/utils"
)

// Dispatcher need for Polling, and webhook
// For Bot run,
// Bot struct uses as API wrapper
// Dispatcher uses as Bot starter
// Another level of abstraction
type Dispatcher struct {
	Bot      *Bot
	Storage  storage.Storage
	Mutex    *sync.Mutex
	Welcome_ bool

	// Handlers
	MessageHandler       HandlerObj
	CallbackQueryHandler HandlerObj
	ChannelPostHandler   HandlerObj
	PollHandler          HandlerObj
	ChatMemberHandler    HandlerObj
	PollAnswerHandler    HandlerObj
	MyChatMemberHandler  HandlerObj

	// If you want to add onshutdown function
	// just append to this object, :P
	OnShutdownCallbacks []*OnStartAndShutdownFunc
	OnStartupCallbacks  []*OnStartAndShutdownFunc

	// private fields
	updatesCh     chan *objects.Update
	currentUpdate *objects.Update
	synchronus    bool
	polling       bool
	webhook       bool
}

var (
	ErrorTypeAssertion = errors.New("impossible to do type assertion to this callback")
)

type OnStartAndShutdownFunc func(dp *Dispatcher)

// Config for start polling method
// idk where to put this config, configs or dispatcher?
type StartPollingConfig struct {
	GetUpdatesConfig
	Relax        time.Duration
	ResetWebhook bool
	ErrorSleep   uint
	SkipUpdates  bool
	SafeExit     bool
	Timeout      time.Duration
}

func NewStartPollingConf(skip_updates bool) *StartPollingConfig {
	return &StartPollingConfig{
		GetUpdatesConfig: GetUpdatesConfig{
			Timeout: 20,
			Limit:   0,
		},
		Relax:        1 * time.Second,
		ResetWebhook: false,
		ErrorSleep:   5,
		SkipUpdates:  skip_updates,
		SafeExit:     true,
		Timeout:      5 * time.Second,
	}
}

type StartWebhookConfig struct {
	*SetWebhookConfig
	BotURL             string
	Address            string
	Handler            http.Handler
	KeyFile            interface{}
	DropPendingUpdates bool
	SafeExit           bool
}

func NewStartWebhookConf(url string, address string) *StartWebhookConfig {
	return &StartWebhookConfig{
		BotURL:  url,
		Address: address,
	}
}

// NewDispathcer get a new Dispatcher
// And with autoconfiguration, need to run once
func NewDispatcher(bot *Bot, storage storage.Storage, synchronus bool) *Dispatcher {
	dp := &Dispatcher{
		Bot:        bot,
		synchronus: synchronus,
		Storage:    storage,
		updatesCh:  make(chan *objects.Update, 1),
	}

	dp.MessageHandler = NewDHandlerObj(dp)
	dp.CallbackQueryHandler = NewDHandlerObj(dp)
	dp.ChannelPostHandler = NewDHandlerObj(dp)
	dp.ChatMemberHandler = NewDHandlerObj(dp)
	dp.PollHandler = NewDHandlerObj(dp)
	dp.PollAnswerHandler = NewDHandlerObj(dp)
	dp.ChannelPostHandler = NewDHandlerObj(dp)

	return dp
}

func (dp *Dispatcher) SetState(state *fsm.State) {
	cid, uid := utils.GetUidAndCidFromUpd(dp.currentUpdate)
	dp.Storage.SetState(cid, uid, state.GetFullState())
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

// RegisterMessageHandler excepts you pass to parametrs a your function
func (dp *Dispatcher) RegisterMessageHandler(callback HandlerFunc) {
	dp.MessageHandler.Register(callback)
}

// ProcessOneUpdate you guess, processes ONLY one comming update
// Support only one Message update
func (dp *Dispatcher) ProcessOneUpdate(update *objects.Update) error {
	var err error

	// very bad code, please dont see this bullshit
	// ============================================
	if update.Message != nil {
		dp.MessageHandler.TriggerMiddleware(update, PREMIDDLEWARE)
		for _, h := range dp.MessageHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.Message))
			if !ok {
				return errors.New("Message handler type assertion error, need type func(*Message), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}

			err = dp.MessageHandler.TriggerMiddleware(update, PROCESSMIDDLEWARE)
			if err != nil {
				log.Println(err)
				continue
			}

			h.Call(update, func() { cb(update.Message) }, dp.synchronus)
		}
		dp.MessageHandler.TriggerMiddleware(update, POSTMIDDLEWARE)

	} else if update.CallbackQuery != nil {
		dp.CallbackQueryHandler.TriggerMiddleware(update, PREMIDDLEWARE)
		for _, h := range dp.CallbackQueryHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.CallbackQuery))
			if !ok {
				return errors.New("Callbackquery handler type assertion error, need type func(*CallbackQuery), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.CallbackQueryHandler.TriggerMiddleware(update, PROCESSMIDDLEWARE)
			if err != nil {
				log.Println(err)
				continue
			}

			h.Call(update, func() { cb(update.CallbackQuery) }, dp.synchronus)
		}
		dp.CallbackQueryHandler.TriggerMiddleware(update, POSTMIDDLEWARE)

	} else if update.ChannelPost != nil {
		dp.ChannelPostHandler.TriggerMiddleware(update, PREMIDDLEWARE)
		for _, h := range dp.ChannelPostHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.Message))
			if !ok {
				return errors.New("ChannelPost handler type assertion error, need type func(*ChannelPost), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.ChannelPostHandler.TriggerMiddleware(update, PROCESSMIDDLEWARE)
			if err != nil {
				log.Println(err)
				continue
			}

			h.Call(update, func() { cb(update.ChannelPost) }, dp.synchronus)
		}
		dp.ChannelPostHandler.TriggerMiddleware(update, POSTMIDDLEWARE)

	} else if update.Poll != nil {
		dp.PollHandler.TriggerMiddleware(update, PREMIDDLEWARE)
		for _, h := range dp.PollHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.Poll))
			if !ok {
				return errors.New("Poll handler type assertion error, need type func(*Poll), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.PollHandler.TriggerMiddleware(update, PROCESSMIDDLEWARE)
			if err != nil {
				log.Println(err)
				continue
			}

			h.Call(update, func() { cb(update.Poll) }, dp.synchronus)
		}
		dp.PollHandler.TriggerMiddleware(update, POSTMIDDLEWARE)

	} else if update.PollAnswer != nil {
		dp.PollAnswerHandler.TriggerMiddleware(update, PREMIDDLEWARE)
		for _, h := range dp.PollAnswerHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.PollAnswer))
			if !ok {
				return errors.New("PollAnswer handler type assertion error, need type func(*PollAnswer), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.PollAnswerHandler.TriggerMiddleware(update, PROCESSMIDDLEWARE)
			if err != nil {
				log.Println(err)
				continue
			}

			h.Call(update, func() { cb(update.PollAnswer) }, dp.synchronus)
		}
		dp.PollAnswerHandler.TriggerMiddleware(update, POSTMIDDLEWARE)

	} else if update.ChatMember != nil {
		dp.ChatMemberHandler.TriggerMiddleware(update, PREMIDDLEWARE)
		for _, h := range dp.ChatMemberHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.ChatMember))
			if !ok {
				return errors.New("ChatMember handler type assertion error, need type func(*ChatMember), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.ChatMemberHandler.TriggerMiddleware(update, PROCESSMIDDLEWARE)
			if err != nil {
				log.Println(err)
				continue
			}

			h.Call(update, func() { cb(update.ChatMember) }, dp.synchronus)
		}
		dp.ChatMemberHandler.TriggerMiddleware(update, POSTMIDDLEWARE)

	} else if update.MyChatMember != nil {
		dp.MyChatMemberHandler.TriggerMiddleware(update, PREMIDDLEWARE)
		for _, h := range dp.MyChatMemberHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.ChatMemberUpdated))
			if !ok {
				return errors.New("MyChatMember handler type assertion error, need type func(*ChatMemberUpdated), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.MyChatMemberHandler.TriggerMiddleware(update, PROCESSMIDDLEWARE)
			if err != nil {
				log.Println(err)
				continue
			}

			h.Call(update, func() { cb(update.MyChatMember) }, dp.synchronus)
		}
		dp.MyChatMemberHandler.TriggerMiddleware(update, POSTMIDDLEWARE)

	} else {
		text := "detected not supported type of updates seems like telegram bot api updated before this package updated"
		return errors.New(text)
	}

	// end of adventure
	return nil
}

// SkipUpdates skip comming updates, sending to telegram servers
func (dp *Dispatcher) SkipUpdates() {
	dp.Bot.GetUpdates(&GetUpdatesConfig{
		Offset:  -1,
		Timeout: 1,
	})
}

// ========================================
// On Startup and Shutdown related methods
// ========================================

// Shutdown calls when you enter ^C(which means SIGINT)
// And SafeExit trap it, before you exit
func (dp *Dispatcher) Shutdown() {
	for _, cb := range dp.OnShutdownCallbacks {
		c := *cb
		if dp.synchronus {
			c(dp)
		} else {
			go c(dp)
		}
	}
}

// StartUp function, iterate over a callbacks from OnStartupCallbacks
// Calls in StartPolling function
func (dp *Dispatcher) StartUp() {
	for _, cb := range dp.OnStartupCallbacks {
		c := *cb
		if dp.synchronus {
			c(dp)
		} else {
			go c(dp)
		}
	}
}

// Onstartup method append to OnStartupCallbaks a callbacks
// Using pointers bc cant unregister function using copy of object
// And golang doesnot support generics, and type equals
func (dp *Dispatcher) OnStartup(f ...OnStartAndShutdownFunc) {
	var objs []*OnStartAndShutdownFunc

	for _, cb := range f {
		objs = append(objs, &cb)
	}

	dp.OnStartupCallbacks = append(dp.OnStartupCallbacks, objs...)
}

// OnShutdown method using for register OnShutdown callbacks
// Same code like OnStartup
func (dp *Dispatcher) OnShutdown(f ...OnStartAndShutdownFunc) {
	var objs []*OnStartAndShutdownFunc

	for _, cb := range f {
		objs = append(objs, &cb)
	}

	dp.OnShutdownCallbacks = append(dp.OnShutdownCallbacks, objs...)
}

// Thanks: https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
func (dp *Dispatcher) SafeExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		dp.ShutDownDP()
		os.Exit(0)
	}()
}

// ShutDownDP calls ResetWebhook for reset webhook in telegram servers, if yes
func (dp *Dispatcher) ShutDownDP() {
	log.Println("Stop polling!")
	dp.ResetWebhook(true)
	dp.Storage.Clear()
	dp.Shutdown()
}

func (dp *Dispatcher) Welcome() {
	dp.Bot.GetMe()
	log.Println("Bot: ", dp.Bot.Me)
}

// GetUpdatesChan makes getUpdates request to telegram servers
// sends update to updates channel
// Time.Sleep here for stop goroutine for a c.Relax time
//
// yeah it bad, and works only on crutches, but works, idk how
func (dp *Dispatcher) MakeUpdatesChan(c *StartPollingConfig) {
	go func() {
		for {
			if c.Relax != 0 {
				time.Sleep(c.Relax)
			}

			updates, err := dp.Bot.GetUpdates(&c.GetUpdatesConfig)
			if err != nil {
				log.Println(err.Error())
				log.Println("Error with getting updates")
				time.Sleep(time.Duration(c.ErrorSleep))

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= c.Offset {
					c.Offset = update.UpdateID + 1
					dp.updatesCh <- update
				}
			}
		}
	}()
}

// StartPolling check out to comming updates
// If yes, Telegram Get to your bot a Update
// Using GetUpdates method in Bot structure
// GetUpdates config using for getUpdates method
func (dp *Dispatcher) StartPolling(c *StartPollingConfig) error {
	if c.SafeExit {
		// runs goroutine for safly terminate program(bot)
		dp.SafeExit()
	}

	dp.StartUp()
	if c.ResetWebhook {
		dp.ResetWebhook(true)
	}
	if c.SkipUpdates {
		dp.SkipUpdates()
	}
	if dp.Welcome_ {
		dp.Welcome()
	}
	// TODO: timeout
	dp.polling = true
	dp.MakeUpdatesChan(c)

	for upd := range dp.updatesCh {
		dp.currentUpdate = upd
		err := dp.ProcessOneUpdate(upd)
		if err != nil {
			return err
		}
	}

	return errors.New(":P")
}

func (dp *Dispatcher) MakeWebhookChan(c *StartWebhookConfig) {
	http.HandleFunc(c.BotURL, func(wr http.ResponseWriter, req *http.Request) {
		update, err := utils.RequestToUpdate(req)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			wr.WriteHeader(http.StatusBadRequest)
			wr.Header().Set("Content-Type", "application/json")
			wr.Write(errMsg)
			return
		}

		dp.updatesCh <- update
	})
}

// StartWebhook method registers on BotUrl uri a function which handles every comming update
// Using In Pair of SetWebhook method
// Startup method executes after SetWebhook method
// and adds a OnShutdown function, and no need to delete webhook with hands
func (dp *Dispatcher) StartWebhook(c *StartWebhookConfig) error {
	if dp.polling {
		return errors.New(
			"you want to enable two conflict modes at the same time, polling and webhook",
		)
	}
	dp.OnShutdown(func(dp *Dispatcher) {
		dw := &DeleteWebhookConfig{
			DropPendingUpdates: c.DropPendingUpdates,
		}
		dp.Bot.DeleteWebhook(dw)
	})
	_, err := dp.Bot.SetWebhook(c.SetWebhookConfig)
	if err != nil {
		return err
	}
	dp.StartUp()
	if c.SafeExit {
		dp.SafeExit()
	}
	http.HandleFunc(c.BotURL, func(wr http.ResponseWriter, req *http.Request) {
		update, err := utils.RequestToUpdate(req)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			wr.WriteHeader(http.StatusBadRequest)
			wr.Header().Set("Content-Type", "application/json")
			wr.Write(errMsg)
			return
		}

		dp.currentUpdate = update
		err = dp.ProcessOneUpdate(update)
		if err != nil {
			fmt.Println(err)
		}
	})
	certPath, err := utils.GuessFileName(c.Certificate)
	if err != nil {
		return err
	}
	keyfile, err := utils.GuessFileName(c.KeyFile)
	if err != nil {
		return err
	}
	dp.webhook = true
	return http.ListenAndServeTLS(c.Address, certPath, keyfile, c.Handler)
}
