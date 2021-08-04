package tgp

import (
	"sync"

	"github.com/pikoUsername/tgp/objects"
)

type Cluster struct {
	dp   *Dispatcher
	rmu  sync.RWMutex
	bots []*Bot
}

func NewCluster(dp *Dispatcher, max_connections int) *Cluster {
	bots := make([]*Bot, max_connections)
	
	return &Cluster{
		dp: dp, 
		bots: , 
	}
}

func (c *Cluster) StartPolling(cpc *StartPollingConfig) error {
	return nil
}

func (c *Cluster) StartWebhook(cwc *StartWebhookConfig) error {
	return nil
}

func (c *Cluster) ProcessUpdates(ch <-chan *objects.Update) error {
	return nil
}
