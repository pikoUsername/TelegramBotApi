simple Telegram Bot Api wrapper, maybe will grow to framework 

<h1>
This package is on alpha version,
and any update can broke backward capability
</h1>

NOTE: Please, don't try to use this package in serious projects 

## version 0.1.1
<br>

## docs
<br>
 (WIP) for first time, you can read the code 

Download This package using `go get -v github.com/pikoUsername/tgp` command 

This package based/expired/copy_pasted on go-telegram-bot, and aiogram

## Example
```go
package main

import (
	"fmt"
	"log"
    
    "github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/utils"
    "github.com/pikoUsername/tgp/objects"
)

func handler(bot *tgp.Bot, m *objects.Message) { 
    if m.Text == "" { 
        return
    }

    _, err := bot_.SendMessage(&tgp.SendMessageConfig{
        ChatID: m.Chat.ID, 
        Text: m.Text, 
    })
    if err != nil { 
        fmt.Println(err)
    }
}

func main() {
	bot, err := tgp.NewBot("<token>", true, utils.ModeHTML)
	if err != nil {
		panic(err)
	}

	dp := tgp.NewDispatcher(bot)
	if err != nil {
		panic(err)
	}
    dp.MessageHandler.Register(handler)
    dp.StartPolling(tgp.NewStartPollingConfig(true))
}
```
