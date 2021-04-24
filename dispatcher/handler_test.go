package dispatcher_test

import (
	"testing"

	"github.com/pikoUsername/tgp/configs"
	"github.com/pikoUsername/tgp/dispatcher"
)

var (
	TestChatID = -534916942
)

func TestRegister(t *testing.T) {
	dp, err := GetDispatcher(t)
	if err != nil {
		t.Error(err)
	}
	dp.ResetWebhook(false)
	dp.MessageHandler.Register(func(mes *dispatcher.Context) interface{} {
		msg, err := mes.Send(&configs.SendMessageConfig{
			ChatID: int64(TestChatID),
		})
		if err != nil {
			panic(err)
		}
		return msg
	})
}
