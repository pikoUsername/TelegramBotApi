package dispatcher

import (
	"errors"

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
// NOTE{
//    Using NewDispatcher, Dispatcher confgures itself
//    But} if you want use own realitzation of dispatcher
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
func (dp *Dispatcher) ProcessUpdates(updates []objects.Update) error {
	return nil // TODO
}

// Dont ask me why i m using this function
func stub(a ...interface{}) {}

// ProcessOneUpdate you guess, processes ONLY one comming update
func (dp *Dispatcher) ProcessOneUpdate(update objects.Update) error {
	var event interface{}
	var update_type string
	// oh shit, why compiler so struct, and when you using if else if and etc.
	// raises not using variable, ;(
	stub(event, update_type)

	if update.Message != nil {
		update_type = "message"
		event = update.Message
	} else if update.EditedMessage != nil {
		update_type = "edited_message"
		event = update.EditedMessage
	} else if update.ChannelPost != nil {
		update_type = "channel_post"
		event = update.ChannelPost
	} else if update.EditedChannelPost != nil {
		update_type = "edited_channel_post"
		event = update.EditedChannelPost
	} else if update.InlineQuery != nil {
		update_type = "inline_query"
		event = update.InlineQuery
	} else if update.ChosenInlineResult != nil {
		update_type = "chosen_inline_result"
		event = update.ChosenInlineResult
	} else if update.CallbackQuery != nil {
		update_type = "callback_query"
		event = update.CallbackQuery
	} else if update.ShippingQuery != nil {
		update_type = "shipping_query"
		event = update.ShippingQuery
		// } else if update.pre_checkout_query != nil {
		// 	update_type = "pre_checkout_query"
		// 	event = update.Pre
	} else if update.Poll != nil {
		update_type = "poll"
		event = update.Poll
	} else if update.PollAnswer != nil {
		update_type = "poll_answer"
		event = update.PollAnswer
	} else if update.MyChatMember != nil {
		update_type = "my_chat_member"
		event = update.MyChatMember
	} else if update.ChatMember != nil {
		update_type = "chat_member"
		event = update.ChatMember
	} else {
		text := "Detected Not supported type of updates Seems like Telegram bot api updated brfore this package updated"
		return errors.New(text)
	}

	return nil
}

// StartPolling check out to comming updates
// If yes, Telegram Get to your bot a Update
// Using GetUpdates method in Bot structure
// GetUpdates config using for getUpdates method
func (dp *Dispatcher) StartPolling(c *configs.GetUpdatesConfig) error {
	dp.Polling = true
	for dp.Polling {
		// TODO{ timeout
		updates, err := dp.Bot.GetUpdates(c)
		if err != nil {
			return err
		}
		if updates != nil {
			// I cant understand how it s works, and where need to use goroutines
			go dp.ProcessPollingUpdates(updates)
		}
	}
	return nil
}

// StartWebhook ...
func (dp *Dispatcher) StartWebhook() error {
	return nil
}
