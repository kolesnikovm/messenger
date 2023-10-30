package kafka

import (
	"github.com/IBM/sarama"
	"github.com/kolesnikovm/messenger/notifier/kafka"
	"github.com/kolesnikovm/messenger/store"
	"github.com/rs/zerolog/log"
)

type Consumer struct {
	ready             chan bool
	MessageAggregator store.Aggregator
}

func (с *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(с.ready)
	return nil
}

func (с *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	const op = "Consumer.ConsumeClaim"

	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Info().Msgf("%s: message channel was closed", op)
				return nil
			}

			entityMessage, err := kafka.ParseMessage(msg.Value)
			if err != nil {
				log.Error().Err(err).Msg("")
				continue
			}

			c.MessageAggregator.Add(entityMessage)

			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
