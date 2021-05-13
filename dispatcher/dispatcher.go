package dispatcher

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/pikoUsername/tgp/bot"
	"github.com/pikoUsername/tgp/configs"
	"github.com/pikoUsername/tgp/objects"
)

// Dispatcher need for Polling, and webhook
// For Bot run,
// Bot struct uses as API wrapper
// Dispatcher uses as Bot starter
// Middlewares uses function
// Another level of abstraction
type Dispatcher struct {
	Bot *bot.Bot

	// Handlers
	MessageHandler       HandlerObj
	CallbackQueryHandler HandlerObj
	ChannelPost          HandlerObj

	// If you want to add onshutdown function
	// just append to this object, :P
	OnShutdownCallbacks []*OnStartAndShutdownFunc
	OnStartupCallbacks  []*OnStartAndShutdownFunc
}

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
	}
}

// NewDispathcer get a new Dispatcher
// And with autoconfiguration, need to run once
func NewDispatcher(bot *bot.Bot) *Dispatcher {
	dp := &Dispatcher{
		Bot: bot,
	}

	dp.MessageHandler = NewDHandlerObj(dp)
	dp.CallbackQueryHandler = NewDHandlerObj(dp)
	dp.ChannelPost = NewDHandlerObj(dp)

	return dp
}

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

// ProcessUpdates using for process updates from any way
func (dp *Dispatcher) ProcessUpdates(updates []*objects.Update, conf *StartPollingConfig) error {
	for _, upd := range updates {
		err := dp.ProcessOneUpdate(upd)
		if err != nil {
			return err
		}
	}
	return nil
}

// ProcessOneUpdate you guess, processes ONLY one comming update
// Support only one Message update
func (dp *Dispatcher) ProcessOneUpdate(update *objects.Update) error {
	if update.Message != nil {
		dp.MessageHandler.Notify(update)
	} else if update.CallbackQuery != nil {
		dp.CallbackQueryHandler.Notify(update)
	} else if update.ChannelPost != nil {
		dp.ChannelPost.Notify(update)
	} else {
		text := "detected not supported type of updates seems like telegram bot api updated before this package updated"
		return errors.New(text)
	}
	return nil
}

func (dp *Dispatcher) SkipUpdates() {
	dp.Bot.GetUpdates(&configs.GetUpdatesConfig{
		Offset:  -1,
		Timeout: 1,
	})
}

// ========================================
// On Startup and Shutdown related methods
// ========================================

func (dp *Dispatcher) Shutdown() {
	for _, cb := range dp.OnShutdownCallbacks {
		c := *cb
		c(dp)
	}
}

func (dp *Dispatcher) StartUp() {
	for _, cb := range dp.OnStartupCallbacks {
		c := *cb
		c(dp)
	}
}

func (dp *Dispatcher) OnStartup(f ...OnStartAndShutdownFunc) {
	var objs []*OnStartAndShutdownFunc

	for _, cb := range f {
		objs = append(objs, &cb)
	}

	dp.OnStartupCallbacks = append(dp.OnStartupCallbacks, objs...)
}

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
		for c := range c {
			if c.String() != "" {
				dp.ShutDownDP()
				os.Exit(0)
			}
		}
	}()
}

// ShutDownDP calls ResetWebhook for reset webhook in telegram servers, if yes
func (dp *Dispatcher) ShutDownDP() {
	// Here will be other methods
	dp.ResetWebhook(true)
	dp.Shutdown()
}

// StartPolling check out to comming updates
// If yes, Telegram Get to your bot a Update
// Using GetUpdates method in Bot structure
// GetUpdates config using for getUpdates method
func (dp *Dispatcher) StartPolling(c *StartPollingConfig) error {
	if c.ResetWebhook {
		dp.ResetWebhook(true)
	}

	if c.SkipUpdates {
		dp.SkipUpdates()
	}

	dp.StartUp()
	for {
		if c.SafeExit {
			dp.SafeExit()
		}

		// TODO: timeout
		updates, err := dp.Bot.GetUpdates(&c.GetUpdatesConfig)
		if err != nil {
			log.Println(err)
			log.Println("Error with getting updates")
			time.Sleep(time.Duration(c.ErrorSleep))

			continue
		}

		if len(updates) > 0 && updates != nil {
			index := len(updates) - 1

			c.Offset = updates[index].UpdateID + 1

			err := dp.ProcessUpdates(updates, c)
			if err != nil {
				return err
			}
		}

		if c.Relax != 0 {
			time.Sleep(c.Relax)
		}
	}
}
