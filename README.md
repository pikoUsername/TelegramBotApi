simple Telegram Bot Api wrapper, maybe will grow to framework 

<h1>
This package is on alpha version,
and any update can broke backward capability
</h1>

NOTE: Please dont try use this package 

## version 0.1.1
<br>

## docs
<br>
 (WIP) for first time, you can read the code 

Download This package using this command `go get -v github.com/pikoUsername/tgp` 
Ignore a warning about no files in root directory, if you want download a new version of this package 
Then delete a past version in folder where saved on.

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

func main() {
	bot, err := tgp.NewBot("<token>", true, utils.ModeHTML)
	if err != nil {
		panic(err)
	}

	dp := tgp.NewDispatcher(bot)
	if err != nil {
		panic(err)
	}
    dp.MessageHandler.Register(func(m *objects.Message) { 
        if m.Test == "" { 
            return
        }

        _, err := bot.SendMessage(&tgp.SendMessageConfig{
            ChatID: m.Chat.ID, 
            Text: m.Text, 
        })
        if err != nil { 
            fmt.Println(err)
        }
    })
    dp.StartPolling(tgp.NewStartPollingConfig(true))
}
```
