package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

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
		MessageID:   msg.MessageID,
		SenderID:    msg.SenderID,
		RecipientID: msg.RecipientID,
		Text:        msg.Text,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	messageKey := composeKey(msg.SenderID, msg.RecipientID)

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

func composeKey(id1, id2 uint64) string {
	slice := []uint64{id1, id2}
	sort.SliceStable(slice, func(i, j int) bool { return slice[i] < slice[j] })

	var stringSlice []string
	for _, id := range slice {
		stringSlice = append(stringSlice, strconv.Itoa(int(id)))
	}

	return strings.Join(stringSlice, ":")
}
