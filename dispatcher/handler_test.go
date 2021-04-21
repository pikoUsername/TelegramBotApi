package dispatcher_test

import (
	"testing"
)

func TestRegister(t *testing.T) {
	dp, err := GetDispatcher(t)
	if err != nil {
		t.Error(err)
	}
	dp.ResetWebhook(false)
	// oh shit, nooo, no way, i cant use objects.Message type,
	// but first argument is interface, is not fair ;(
	// dp.MessageHandler.Register(func(mes interface{} , bot bot.Bot) {

	// })
}
