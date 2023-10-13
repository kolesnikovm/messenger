package kafka

import (
	"github.com/IBM/sarama"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/rs/zerolog/log"
)

type KafkaMessageSender struct {
	Producer sarama.SyncProducer
	Config   configs.KafkaConfig
}

func New(conf configs.KafkaConfig) *KafkaMessageSender {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(conf.BrokerList, config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start kafka producer")
	}

	return &KafkaMessageSender{
		Producer: producer,
		Config:   conf,
	}
}
