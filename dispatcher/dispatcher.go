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
type Dispatcher struct {
	Bot *bot.Bot

	// Problem fixed ;)
	MessageHandler []*func(*objects.Message)
}

// RegisterMessageHandler except func(*objects.Message) but strict typing
func (dp *Dispatcher) RegisterMessageHandler(callback func(*objects.Message)) {
	dp.MessageHandler = append(dp.MessageHandler, &callback)
}

// ProcessUpdates havenot got any efficient
// if you use webhook and long polling
func (dp *Dispatcher) ProcessUpdates(update *objects.Update) {
	return
}

// StartPolling check out to comming updates
// If yes, Telegram Get to your bot a Update
// Using GetUpdates function in Bot structure
func (dp *Dispatcher) StartPolling(timeout int, limit int) {
	for {
		break
	}
}
