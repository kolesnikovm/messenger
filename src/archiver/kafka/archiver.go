package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/cenkalti/backoff/v4"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store"
	"github.com/rs/zerolog/log"
)

type Archiver struct {
	Client        sarama.ConsumerGroup
	GroupConsumer *Consumer
}

const (
	consumerGroup = "archiver"
	messageTopic  = "messages"
)

func New(kafkaConfig configs.Kafka, archiverConfig configs.Archiver, messageStore store.Messages) (*Archiver, error) {
	const op = "Archiver.New"

	consumerConfig := sarama.NewConfig()

	client, err := sarama.NewConsumerGroup(kafkaConfig.BrokerList, consumerGroup, consumerConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	backoff := backoff.NewExponentialBackOff()
	backoff.InitialInterval = archiverConfig.FlushInterval
	backoff.MaxElapsedTime = 0
	backoff.Reset()

	consumer := &Consumer{
		MessageStore: messageStore,
		Config:       archiverConfig,
		Backoff:      backoff,
	}

	return &Archiver{
		Client:        client,
		GroupConsumer: consumer,
	}, nil
}

func (a *Archiver) Start(ctx context.Context) {
	const op = "Archiver.Start"

	go func() {
		for {
			if err := a.Client.Consume(ctx, []string{messageTopic}, a.GroupConsumer); err != nil {
				log.Error().Err(err).Str("op", op).Send()
				return
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()
}

func (a *Archiver) Close() {
	const op = "Archiver.Close"

	err := a.Client.Close()
	if err != nil {
		log.Error().Err(err).Str("op", op).Send()
	}
}
