package tgp

import (
	"fmt"

	"github.com/pikoUsername/tgp/fsm/storage"
)

func createBot(token string, storage storage.Storage) (*Dispatcher, error) {
	bot, err := NewBot(token, "HTML", nil)
	if err != nil {
		return nil, err
	}
	dp := NewDispatcher(bot, storage)
	return dp, nil
}

// RunPolling is more compact version of just dp.RunPolling
func RunPolling(token string, storage storage.Storage) error {
	dp, err := createBot(token, storage)
	if err != nil {
		return err
	}

	return dp.RunPolling(NewPollingConfig(true))
}

// RunWebhook similar to RunPolling function, but webhook
func RunWebhook(token string, address string, storage storage.Storage) error {
	dp, err := createBot(token, storage)
	if err != nil {
		return err
	}
	c := NewWebhookConfig(fmt.Sprintf("/%s", token), address)

	return dp.RunWebhook(c)
}
