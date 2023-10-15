package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/rs/zerolog/log"
)

type kafkaMessage struct {
	Text string `json:"text"`
}

func (k *KafkaMessageSender) Send(ctx context.Context, msg entity.Message) error {
	const op = "KafkaMessageSender.Send"

	payload, err := json.Marshal(&kafkaMessage{
		Text: msg.Text,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	partition, offset, err := k.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: "messages",
		Value: sarama.ByteEncoder(payload),
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
