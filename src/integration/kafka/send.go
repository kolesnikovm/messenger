package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/rs/zerolog/log"
)

func (k *KafkaMessageSender) Send(msg entity.Message) error {
	const op = "KafkaMessageSender.Send"

	partition, offset, err := k.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: k.Config.Topic,
		Value: sarama.StringEncoder(msg.Text),
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug().
		Int32("partition", partition).
		Int64("offset", offset).
		Msg("message sent to kafa")

	return nil
}
