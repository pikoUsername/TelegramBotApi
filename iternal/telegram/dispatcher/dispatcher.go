package dispatcher

import (
	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/bot"
	"github.com/pikoUsername/TelegramBotApiWrapper/iternal/telegram/objects"
)

// Dispatcher need for Polling, and webhook
// For Bot run,
// Bot struct uses as API wrapper
// Dispatcher uses as Bot starter
type Dispatcher struct {
	Bot *bot.Bot
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
