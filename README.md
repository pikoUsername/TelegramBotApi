# TGP 

| golang version 1.15+ 

> Telegram Golang Package

> go get -u github.com/pikoUsername/tgp

# Overview
TGP(bad name) - telegram bot api client library. 
Developed as better version of other same libraries, 
but not tested, made of shit and sticks, and not optimized version. 

# Examples
Here is simple example of echo bot: 
```go
package main

import (
	"fmt"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/fsm/storage"
)

func main() {
	bot, err := tgp.NewBot("<token>", "HTML", nil)

	if err != nil {
		panic(err)
	}

	dp := tgp.NewDispatcher(bot, storage.NewMemoryStorage())

	dp.MessageHandler.Register(func(ctx *tgp.Context) {
		_, err := bot.Send(tgp.NewSendMessage(ctx.Message.Chat.ID, ctx.Message.Text))
		if err != nil {
			fmt.Println(err)
		}
	})
	dp.StartPolling(tgp.NewStartPollingConf(true))
}
``` 
(See more complicated examples in _examples/ directory)

# FSM
Finate State Machine - Sets state for every user in bot,
and uses special filter for interact with user, using Finate States
Currently FSM on testing period, and not well tested
See example in _examples directory. 

# Filters
Filters support - as i mentioned before, filters uses as filters. 
Currently supported Regex, FSM, command, test, content-type, chat-type filters.  

# Warning
If you wanna use this package, for something more than hello, world bot.
Don't do it! Package have no stable version in this moment

Current Version of TGP - 0.13.0v
