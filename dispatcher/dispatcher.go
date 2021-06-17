package dispatcher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/pikoUsername/tgp/bot"
	"github.com/pikoUsername/tgp/configs"
	"github.com/pikoUsername/tgp/dispatcher/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

// Dispatcher need for Polling, and webhook
// For Bot run,
// Bot struct uses as API wrapper
// Dispatcher uses as Bot starter
// Middlewares uses function
// Another level of abstraction
type Dispatcher struct {
	Bot     *bot.Bot
	Storage storage.Storage

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

	synchronus bool
}

var (
	ErrorTypeAssertion = errors.New("can not do type assertion to this callback")
)

type OnStartAndShutdownFunc func(dp *Dispatcher)

// Config for start polling method
// idk where to put this config, configs or dispatcher?
type StartPollingConfig struct {
	configs.GetUpdatesConfig
	Relax        time.Duration
	ResetWebhook bool
	ErrorSleep   uint
	SkipUpdates  bool
	SafeExit     bool
	Timeout      time.Duration
}

func NewStartPollingConf(skip_updates bool) *StartPollingConfig {
	return &StartPollingConfig{
		GetUpdatesConfig: configs.GetUpdatesConfig{
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

// NewDispathcer get a new Dispatcher
// And with autoconfiguration, need to run once
func NewDispatcher(bot *bot.Bot, storage storage.Storage, synchronus bool) *Dispatcher {
	dp := &Dispatcher{
		Bot:        bot,
		synchronus: synchronus,
		Storage:    storage,
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

// ResetWebhook uses for reset webhook for telegram
func (dp *Dispatcher) ResetWebhook(check bool) error {
	if check {
		wi, err := dp.Bot.GetWebhookInfo()
		if err != nil {
			return err
		}
		if wi.URL == "" {
			return nil
		}
	}
	return dp.Bot.DeleteWebhook(&configs.DeleteWebhookConfig{})
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
		for _, h := range dp.MessageHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.Message))
			if !ok {
				return errors.New("Message handler type assertion error, need type func(*Message), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}

			err = dp.MessageHandler.TriggerMiddleware(update)
			if err != nil {
				return err
			}

			h.Call(update, func() { cb(update.Message) }, dp.synchronus)
		}

	} else if update.CallbackQuery != nil {
		for _, h := range dp.CallbackQueryHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.CallbackQuery))
			if !ok {
				return errors.New("Callbackquery handler type assertion error, need type func(*CallbackQuery), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.CallbackQueryHandler.TriggerMiddleware(update)
			if err != nil {
				return err
			}

			h.Call(update, func() { cb(update.CallbackQuery) }, dp.synchronus)
		}

	} else if update.ChannelPost != nil {
		for _, h := range dp.ChannelPostHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.Message))
			if !ok {
				return errors.New("ChannelPost handler type assertion error, need type func(*ChannelPost), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.ChannelPostHandler.TriggerMiddleware(update)
			if err != nil {
				return err
			}

			h.Call(update, func() { cb(update.ChannelPost) }, dp.synchronus)
		}

	} else if update.Poll != nil {
		for _, h := range dp.PollHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.Poll))
			if !ok {
				return errors.New("Poll handler type assertion error, need type func(*Poll), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.PollHandler.TriggerMiddleware(update)
			if err != nil {
				return err
			}

			h.Call(update, func() { cb(update.Poll) }, dp.synchronus)
		}

	} else if update.PollAnswer != nil {
		for _, h := range dp.PollAnswerHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.PollAnswer))
			if !ok {
				return errors.New("PollAnswer handler type assertion error, need type func(*PollAnswer), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.PollAnswerHandler.TriggerMiddleware(update)
			if err != nil {
				return err
			}

			h.Call(update, func() { cb(update.PollAnswer) }, dp.synchronus)
		}

	} else if update.ChatMember != nil {
		for _, h := range dp.ChatMemberHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.ChatMember))
			if !ok {
				return errors.New("ChatMember handler type assertion error, need type func(*ChatMember), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.ChatMemberHandler.TriggerMiddleware(update)
			if err != nil {
				return err
			}

			h.Call(update, func() { cb(update.ChatMember) }, dp.synchronus)
		}

	} else if update.MyChatMember != nil {
		for _, h := range dp.MyChatMemberHandler.GetHandlers() {
			i_cb := *h.Callback
			cb, ok := i_cb.(func(*objects.ChatMemberUpdated))
			if !ok {
				return errors.New("MyChatMember handler type assertion error, need type func(*ChatMemberUpdated), current type is - " + fmt.Sprintln(reflect.TypeOf(i_cb)))
			}
			err = dp.MyChatMemberHandler.TriggerMiddleware(update)
			if err != nil {
				return err
			}

			h.Call(update, func() { cb(update.MyChatMember) }, dp.synchronus)
		}

	} else {
		text := "detected not supported type of updates seems like telegram bot api updated before this package updated"
		return errors.New(text)
	}

	// end of something
	return nil
}

// SkipUpdates skip comming updates, sending to telegram servers
func (dp *Dispatcher) SkipUpdates() {
	go dp.Bot.GetUpdates(&configs.GetUpdatesConfig{
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
		c(dp)
	}
}

// StartUp function, iterate over a callbacks from OnStartupCallbacks
// Calls in StartPolling function
func (dp *Dispatcher) StartUp() {
	for _, cb := range dp.OnStartupCallbacks {
		c := *cb
		c(dp)
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
	if dp.synchronus {
		dp.Shutdown()
	} else {
		go dp.Shutdown()
	}
}

// GetUpdatesChan makes getUpdates request to telegram servers
// sends update to updates channel
// Time.Sleep here for stop goroutine for a c.Relax time
//
// yeah it bad, and works only on crutches, but works, idk how
func (dp *Dispatcher) GetUpdatesChan(c *StartPollingConfig) chan *objects.Update {
	upd_c := make(chan *objects.Update, c.Limit)

	go func() {
		for {
			if c.Relax != 0 {
				time.Sleep(c.Relax)
			}

			updates, err := dp.Bot.GetUpdates(&c.GetUpdatesConfig)
			if err != nil {
				log.Println(err)
				log.Println("Error with getting updates")
				time.Sleep(time.Duration(c.ErrorSleep))

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= c.Offset {
					c.Offset = update.UpdateID + 1
					upd_c <- update
				}
			}
		}
	}()

	return upd_c
}

// StartPolling check out to comming updates
// If yes, Telegram Get to your bot a Update
// Using GetUpdates method in Bot structure
// GetUpdates config using for getUpdates method
func (dp *Dispatcher) StartPolling(c *StartPollingConfig) error {
	if c.SafeExit {
		// runs goroutine for safly terminate program(bot)
		go dp.SafeExit()
	}

	dp.StartUp()
	if c.ResetWebhook {
		dp.ResetWebhook(true)
	}
	if c.SkipUpdates {
		dp.SkipUpdates()
	}

	// TODO: timeout
	updates := dp.GetUpdatesChan(c)

	for upd := range updates {
		// waiting untill update come...
		if upd != nil {
			err := dp.ProcessOneUpdate(upd)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
