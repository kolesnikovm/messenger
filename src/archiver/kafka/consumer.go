package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/cenkalti/backoff/v4"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/notifier/kafka"
	"github.com/kolesnikovm/messenger/store"
	"github.com/rs/zerolog/log"
)

type Consumer struct {
	MessageStore store.Messages
	Config       configs.Archiver
	Backoff      *backoff.ExponentialBackOff
}

func (с *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (с *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	const op = "Consumer.ConsumeClaim"

	messages := make([]*entity.Message, 0, c.Config.BatchSize)

	ticker := time.NewTicker(c.Config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Info().Msgf("%s: message channel was closed", op)
				return nil
			}

			entityMessage, err := kafka.ParseMessage(msg.Value)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}

			messages = append(messages, entityMessage)

			if len(messages) == c.Config.BatchSize {
				c.sendMessages(context.Background(), messages)

				session.MarkMessage(msg, "")

				messages = messages[:0]
			}
		case <-ticker.C:
			if len(messages) > 0 {
				c.sendMessages(context.Background(), messages)
			}

			messages = messages[:0]
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *Consumer) sendMessages(ctx context.Context, messages []*entity.Message) {
	operation := func() error {
		return c.MessageStore.BatchInsert(ctx, messages)
	}

	notify := func(err error, delay time.Duration) {
		log.Error().Err(err).Dur("delay", delay).Send()
	}

	err := backoff.RetryNotify(operation, c.Backoff, notify)
	if err == nil {
		c.Backoff.Reset()
		return
	}

	log.Error().Err(err).Msg("failed to save messages to db")
}
