package aggregator

import (
	"context"
	"time"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/store"
	"github.com/rs/zerolog/log"
)

type Aggregator struct {
	MessageStore store.Messages
	Messages     chan *entity.Message
	Config       configs.Store
}

func New(conf configs.Store, messageStore store.Messages) *Aggregator {
	messageChan := make(chan *entity.Message, conf.BatchSize)

	return &Aggregator{
		MessageStore: messageStore,
		Messages:     messageChan,
		Config:       conf,
	}
}

func (a *Aggregator) Start(ctx context.Context) {
	for i := 0; i < a.Config.NumSenders; i++ {
		go a.aggregate(ctx)
	}
}

func (a *Aggregator) aggregate(ctx context.Context) {
	messages := make([]*entity.Message, 0, a.Config.BatchSize)

	ticker := time.NewTicker(a.Config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case msg := <-a.Messages:
			messages = append(messages, msg)

			if len(messages) == a.Config.BatchSize {
				a.sendMessages(context.Background(), messages)

				messages = make([]*entity.Message, 0, a.Config.BatchSize)
			}
		case <-ticker.C:
			a.sendMessages(context.Background(), messages)

			messages = make([]*entity.Message, 0, a.Config.BatchSize)
		case <-ctx.Done():
			return
		}
	}
}

func (a *Aggregator) sendMessages(ctx context.Context, messages []*entity.Message) {
	err := a.MessageStore.BatchInsert(ctx, messages)
	if err == nil {
		return
	}

	ticker := time.NewTicker(a.Config.FlushInterval)
	defer ticker.Stop()

	for err != nil {
		log.Error().Err(err).Msg("")

		select {
		case <-ticker.C:
			err = a.MessageStore.BatchInsert(ctx, messages)
		case <-ctx.Done():
			return
		}
	}
}
