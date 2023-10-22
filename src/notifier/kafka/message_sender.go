package kafka

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type KafkaMessageSender struct {
	Producer           sarama.SyncProducer
	Consumer           sarama.Consumer
	PartitionConsumers map[int32]sarama.PartitionConsumer
	Streams            map[uint64]map[ulid.ULID](chan *entity.Message)
	Config             configs.KafkaConfig
}

type kafkaMessage struct {
	MessageID   string `json:"messageId"`
	SenderID    uint64 `json:"senderId"`
	RecipientID uint64 `json:"recipientId"`
	Text        string `json:"text"`
}

const messageTopic = "messages"

func New(conf configs.KafkaConfig) (*KafkaMessageSender, error) {
	const op = "KafkaMessageSender.New"

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(conf.BrokerList, config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	consumer, err := sarama.NewConsumer(conf.BrokerList, config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	partitions, err := consumer.Partitions(messageTopic)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	partitionConsumers := make(map[int32]sarama.PartitionConsumer)
	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(messageTopic, partition, sarama.OffsetOldest)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		partitionConsumers[partition] = partitionConsumer
	}

	streams := make(map[uint64]map[ulid.ULID](chan *entity.Message))

	kafkaMessageSender := &KafkaMessageSender{
		Producer:           producer,
		Consumer:           consumer,
		PartitionConsumers: partitionConsumers,
		Streams:            streams,
		Config:             conf,
	}

	// TODO handle context
	kafkaMessageSender.startConsumers(context.Background())

	return kafkaMessageSender, nil
}

func (k *KafkaMessageSender) Close() {
	for _, partitionConsumer := range k.PartitionConsumers {
		partitionConsumer.Close()
	}

	k.Consumer.Close()
	k.Producer.Close()

	for _, userStreams := range k.Streams {
		for _, stream := range userStreams {
			_, open := <-stream
			if open {
				// TODO add mutex
				close(stream)
			}
		}
	}
}

func (k *KafkaMessageSender) startConsumers(ctx context.Context) {
	for _, partitionConsumer := range k.PartitionConsumers {
		go func(partitionConsumer sarama.PartitionConsumer) {
			for {
				select {
				case msg, ok := <-partitionConsumer.Messages():
					if !ok {
						return
					}

					recepientID := binary.LittleEndian.Uint64(msg.Key)
					userStreams, ok := k.Streams[recepientID]
					if !ok {
						continue
					}

					kafkaMessage := &kafkaMessage{}
					err := json.Unmarshal(msg.Value, kafkaMessage)
					if err != nil {
						log.Error().Err(err).Msg("failed to unmarshal message")
						return
					}

					messageID, err := ulid.Parse(kafkaMessage.MessageID)
					if err != nil {
						log.Error().Err(err).Msgf("failed to parse message id from: %s", kafkaMessage.MessageID)
						return
					}

					entityMessage := &entity.Message{
						MessageID:   messageID,
						SenderID:    kafkaMessage.SenderID,
						RecipientID: kafkaMessage.RecipientID,
						Text:        kafkaMessage.Text,
					}

					for _, stream := range userStreams {
						stream <- entityMessage
					}
				case <-ctx.Done():
					return
				}
			}

		}(partitionConsumer)
	}
}
