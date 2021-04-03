package dispatcher

import (
	"github.com/pikoUsername/tgp/bot"
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

	MessageHandler *HandlerObj
}

// NewDispathcer get a new Dispatcher
// And with autoconfiguration, need to run once
func NewDispatcher(bot *bot.Bot) (*Dispatcher, error) {
	dp := &Dispatcher{
		Bot: bot,
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

// RegisterHandler excepts you pass to parametrs a your function
// which have no returns
func (dp *Dispatcher) RegisterHandler(callback func(interface{}, *bot.Bot)) {
	dp.MessageHandler.Register(&callback)
}

// ProcessUpdates havenot got any efficient
// if you use webhook and long polling
func (dp *Dispatcher) ProcessUpdates(update *objects.Update) {
	return
}

// StartPolling check out to comming updates
// If yes, Telegram Get to your bot a Update
// Using GetUpdates method in Bot structure
func (dp *Dispatcher) StartPolling(timeout int, limit int) error {
	for {
		break
	}
	return nil
}

// StartWebhook ...
func (dp *Dispatcher) StartWebhook() error {
	return nil
}
