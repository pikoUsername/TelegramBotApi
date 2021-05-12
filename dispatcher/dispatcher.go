package dispatcher

import (
	"errors"
	"log"
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
}

// Config for start polling method
// idk where to put this config, configs or dispatcher?
type StartPollingConfig struct {
	configs.GetUpdatesConfig
	Relax        time.Duration
	ResetWebhook bool
	ErrorSleep   uint
	SkipUpdates  bool
}

func NewStartPollingConf(skip_updates bool) *StartPollingConfig {
	return &StartPollingConfig{
		GetUpdatesConfig: configs.GetUpdatesConfig{
			Timeout: 20,
			Limit:   0,
		},
		Relax:        1,
		ResetWebhook: false,
		ErrorSleep:   5,
		SkipUpdates:  skip_updates,
	}
}

// NewDispathcer get a new Dispatcher
// And with autoconfiguration, need to run once
func NewDispatcher(bot *bot.Bot) (*Dispatcher, error) {
	dp := &Dispatcher{
		Bot: bot,
	}

	dp.MessageHandler = NewDHandlerObj(dp)
	dp.CallbackQueryHandler = NewDHandlerObj(dp)
	dp.ChannelPost = NewDHandlerObj(dp)

	return dp, nil
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
		text := "Detected not supported type of updates seems like Telegram bot api updated brfore this package updated"
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

	for {
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

// StartWebhook ...
func (dp *Dispatcher) StartWebhook() error {
	return nil
}
