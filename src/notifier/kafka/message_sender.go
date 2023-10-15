package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/kolesnikovm/messenger/configs"
)

type KafkaMessageSender struct {
	Producer sarama.SyncProducer
	Config   configs.KafkaConfig
}

func New(conf configs.KafkaConfig) (*KafkaMessageSender, error) {
	const op = "KafkaMessageSender.New"

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(conf.BrokerList, config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &KafkaMessageSender{
		Producer: producer,
		Config:   conf,
	}, nil
}
