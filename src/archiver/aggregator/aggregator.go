package aggregator

import (
	"time"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/store"
)

type Aggregator struct {
	MessageStore store.Messages
	Messages     []*entity.Message
	Ticker       *time.Ticker
	Config       configs.Store
}

func New(conf configs.Store, messageStore store.Messages) *Aggregator {
	messages := make([]*entity.Message, 0, conf.BatchSize)

	ticker := time.NewTicker(conf.FlushInterval)

	return &Aggregator{
		MessageStore: messageStore,
		Messages:     messages,
		Ticker:       ticker,
		Config:       conf,
	}
}
