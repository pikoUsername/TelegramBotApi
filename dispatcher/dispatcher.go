package dispatcher

import (
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
// NOTE:
//    Using NewDispatcher, Dispatcher confgures itself
//    But if you want use own realitzation of dispatcher
//    Use Configure method of Dispatcher
type Dispatcher struct {
	Bot *bot.Bot

	// Handlers
	MessageHandler       *HandlerObj
	CallbackQueryHandler *HandlerObj

	Polling bool
}

// NewDispathcer get a new Dispatcher
// And with autoconfiguration, need to run once
func NewDispatcher(bot *bot.Bot) (*Dispatcher, error) {
	dp := &Dispatcher{
		Bot:     bot,
		Polling: false,
	}
	dp.Configure()
	return dp, nil
}

// Configure method Recreadtes all Handlers
// Be care ful about it, but lost registered handlers
// is not scary
func (dp *Dispatcher) Configure() {
	dp.MessageHandler = &HandlerObj{}
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

// RegisterHandler excepts you pass to parametrs a your function
// which have no returns
func (dp *Dispatcher) RegisterHandler(callback *func(interface{}, bot.Bot)) {
	dp.MessageHandler.Register(callback)
}

// ProcessUpdates havenot got any efficient
// if you use webhook and long polling
func (dp *Dispatcher) ProcessPollingUpdates(updates []objects.Update) error {
	return nil // TODO
}

// ProcessUpdates using for process updates from any way
func (dp *Dispatcher) ProcessUpdates(updates []*objects.Update) error {
	return nil // TODO
}

// StartPolling check out to comming updates
// If yes, Telegram Get to your bot a Update
// Using GetUpdates method in Bot structure
// GetUpdates config using for getUpdates method
func (dp *Dispatcher) StartPolling(c *configs.GetUpdatesConfig) error {
	dp.Polling = true
	for dp.Polling {
		// TODO: timeout
		updates, err := dp.Bot.GetUpdates(c)
		if err != nil {
			return err
		}
		if updates != nil {
			dp.ProcessPollingUpdates(updates)
		}
	}
	return nil
}

// StartWebhook ...
func (dp *Dispatcher) StartWebhook() error {
	return nil
}
