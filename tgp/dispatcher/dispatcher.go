package dispatcher

import (
	"github.com/pikoUsername/tgp/tgp/bot"
	"github.com/pikoUsername/tgp/tgp/objects"
)

// Dispatcher need for Polling, and webhook
// For Bot run,
// Bot struct uses as API wrapper
// Dispatcher uses as Bot starter
// Middlewares uses function
type Dispatcher struct {
	Bot            *bot.Bot
	MessageHandler *HandlerObj
}

// RegisterMessageHandler except func(*objects.Message) but strict typing
// but Golang dont love that, so you PLEASE make your callback first arguemnt *objects.Message
func (dp *Dispatcher) RegisterNessageHandler(callback func(interface{})) {
	append(dp.MessageHandler.Callabacks, callback)
}

// ProcessUpdates havenot got any efficient
// if you use webhook and long polling
func (dp *Dispatcher) ProcessUpdates(update *objects.Update) {
	return
}

// SetWebhook make subscribe to telegram events
// or sends to telegram a request for make
// Subscribe to specific IP, and when user
// sends a message to your bot, Telegram know
// Your bot IP and sends to your bot a Update
// https://core.telegram.org/bots/api#setwebhook
func (dp *Dispatcher) SetWebhook(config *WebhookConfig) error {
	return nil
}

// StartPolling check out to comming updates
// If yes, Telegram Get to your bot a Update
// Using GetUpdates function in Bot structure
func (dp *Dispatcher) StartPolling(timeout int, limit int) {
	for {
		break
	}
}
