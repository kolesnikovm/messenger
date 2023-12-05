package kafka

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/metrics"
	"github.com/kolesnikovm/messenger/notifier/hub"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type KafkaMessageSender struct {
	Producer           sarama.SyncProducer
	Consumer           sarama.Consumer
	PartitionConsumers map[int32]sarama.PartitionConsumer
	StreamHub          *hub.StreamHub
	Config             configs.Kafka
}

type kafkaMessage struct {
	MessageID   ulid.ULID `json:"messageId"`
	SenderID    uint64    `json:"senderId"`
	RecipientID uint64    `json:"recipientId"`
	Text        string    `json:"text"`
}

const messageTopic = "messages"

func New(conf configs.Kafka) (*KafkaMessageSender, error) {
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
		partitionConsumer, err := consumer.ConsumePartition(messageTopic, partition, sarama.OffsetNewest)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		partitionConsumers[partition] = partitionConsumer
	}

	streamHub := hub.New()

	kafkaMessageSender := &KafkaMessageSender{
		Producer:           producer,
		Consumer:           consumer,
		PartitionConsumers: partitionConsumers,
		StreamHub:          streamHub,
		Config:             conf,
	}

	return kafkaMessageSender, nil
}

func (k *KafkaMessageSender) Close() {
	for _, partitionConsumer := range k.PartitionConsumers {
		partitionConsumer.Close()
	}

	k.Consumer.Close()
	k.Producer.Close()
}

func (k *KafkaMessageSender) Start(ctx context.Context) {
	for partition, partitionConsumer := range k.PartitionConsumers {
		go func(partitionConsumer sarama.PartitionConsumer, partition int32) {
			for {
				select {
				case msg, ok := <-partitionConsumer.Messages():
					if !ok {
						log.Error().Int32("partition", partition).Msg("partition consumer channel closed")
						partitionConsumer.AsyncClose()
						return
					}

					recipientIDs, err := getRecipientIDs(msg.Key)
					if err != nil {
						log.Error().Msgf("failed to get recipient ids from kafka key: %s", msg.Key)
						continue
					}

					var streams [][]chan *entity.Message

					for _, recipientID := range recipientIDs {
						userStreams := k.StreamHub.GetStreams(recipientID)
						if len(userStreams) == 0 {
							continue
						}

						streams = append(streams, userStreams)
					}
					if len(streams) == 0 {
						continue
					}

					entityMessage, err := ParseMessage(msg.Value)
					if err != nil {
						log.Error().Err(err).Msg("")
						continue
					}

					for _, header := range msg.Headers {
						if string(header.Key) == "timestamp" {
							sendTime := float64(binary.LittleEndian.Uint64(header.Value))
							receiveTime := float64(time.Now().UnixMilli())

							metrics.MessagesLatency.Observe(receiveTime - sendTime)
						}
					}

					for _, userStreams := range streams {
						for _, stream := range userStreams {
							stream <- entityMessage
						}
					}
				case <-ctx.Done():
					partitionConsumer.AsyncClose()
					return
				}
			}

		}(partitionConsumer, partition)
	}
}

func ParseMessage(byteMessage []byte) (*entity.Message, error) {
	const op = "kafka.ParseMessage"

	kafkaMessage := &kafkaMessage{}

	if err := json.Unmarshal(byteMessage, kafkaMessage); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return entity.NewMessage(
		kafkaMessage.MessageID,
		kafkaMessage.SenderID,
		kafkaMessage.RecipientID,
		kafkaMessage.Text,
	), nil
}

func getRecipientIDs(key []byte) ([]uint64, error) {
	const op = "kafka.getRecipientIDs"

	chatID := string(key)

	switch entity.GetChatType(chatID) {
	case entity.Group:
		_, err := entity.GetGroupID(chatID)
		if err != nil {
			return nil, err
		}

		// TODO return group members
		return nil, nil
	case entity.Channel:
		_, err := entity.GetChannelID(chatID)
		if err != nil {
			return nil, err
		}

		// TODO return channel members
		return nil, nil
	case entity.P2P:
		user1, user2, err := entity.GetUserIDs(chatID)
		if err != nil {
			return nil, err
		}

		return []uint64{user1, user2}, nil
	}

	return nil, fmt.Errorf("%s: failed to determine recipient ids from: %s", op, chatID)
}
