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

type result struct {
	Partition int32
	Offset    int64
	Error     error
}

func (k *KafkaMessageSender) Send(ctx context.Context, msg entity.Message) error {
	const op = "KafkaMessageSender.Send"

	payload, err := json.Marshal(&kafkaMessage{
		MessageID:   msg.MessageID.String(),
		SenderID:    msg.SenderID,
		RecipientID: msg.RecipientID,
		Text:        msg.Text,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	messageKey := make([]byte, 8)
	binary.LittleEndian.PutUint64(messageKey, msg.RecipientID)

	resultCh := make(chan *result)

	go func() {
		defer close(resultCh)

		partition, offset, err := k.Producer.SendMessage(&sarama.ProducerMessage{
			Key:   sarama.ByteEncoder(messageKey),
			Topic: messageTopic,
			Value: sarama.ByteEncoder(payload),
		})

		resultCh <- &result{
			Partition: partition,
			Offset:    offset,
			Error:     err,
		}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, context.Cause(ctx))
	case result := <-resultCh:
		if result.Error != nil {
			return fmt.Errorf("%s: %w", op, result.Error)
		}
		log.Debug().
			Int32("partition", result.Partition).
			Int64("offset", result.Offset).
			Msg("message sent to kafa")
	}

	return nil
}
