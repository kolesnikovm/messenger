package aggregator

import (
	"context"
	"time"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/rs/zerolog/log"
)

func (a *Aggregator) Flush(msg *entity.Message) bool {
	a.Messages = append(a.Messages, msg)

	select {
	case <-a.Ticker.C:
		a.sendMessages(context.Background(), a.Messages)

		a.Messages = make([]*entity.Message, 0, a.Config.BatchSize)

		return true
	default:
		if len(a.Messages) == a.Config.BatchSize {
			a.sendMessages(context.Background(), a.Messages)

			a.Messages = make([]*entity.Message, 0, a.Config.BatchSize)

			return true
		}
	}

	return false
}

func (a *Aggregator) sendMessages(ctx context.Context, messages []*entity.Message) {
	err := a.MessageStore.BatchInsert(ctx, messages)
	if err == nil {
		return
	}

	ticker := time.NewTicker(a.Config.FlushInterval)
	defer ticker.Stop()

	for err != nil {
		log.Error().Err(err).Send()

		select {
		case <-ticker.C:
			err = a.MessageStore.BatchInsert(ctx, messages)
		case <-ctx.Done():
			log.Error().Msg("failed to save messages to db")
			return
		}
	}
}
