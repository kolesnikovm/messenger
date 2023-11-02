package kafka

import (
	"github.com/IBM/sarama"
	"github.com/kolesnikovm/messenger/archiver"
	"github.com/kolesnikovm/messenger/notifier/kafka"
	"github.com/rs/zerolog/log"
)

type Consumer struct {
	ready             chan bool
	MessageAggregator archiver.Aggregator
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

	messages := make([]*sarama.ConsumerMessage, 0)

	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Info().Msgf("%s: message channel was closed", op)
				return nil
			}

			messages = append(messages, msg)

			entityMessage, err := kafka.ParseMessage(msg.Value)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}

			flushed := c.MessageAggregator.Flush(entityMessage)

			if flushed {
				for _, message := range messages {
					session.MarkMessage(message, "")
				}

				messages = make([]*sarama.ConsumerMessage, 0)
			}
		case <-session.Context().Done():
			return nil
		}
	}
}
