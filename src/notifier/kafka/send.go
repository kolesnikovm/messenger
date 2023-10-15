package kafka

import (
	"context"
	"encoding/binary"
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

	messageKey := make([]byte, 8)
	binary.LittleEndian.PutUint64(messageKey, msg.RecipientID)

	partition, offset, err := k.Producer.SendMessage(&sarama.ProducerMessage{
		Key:   sarama.ByteEncoder(messageKey),
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
