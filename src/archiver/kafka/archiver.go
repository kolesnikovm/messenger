package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
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

func New(ctx context.Context, conf configs.KafkaConfig, messageAggregator store.Aggregator) (*Archiver, error) {
	const op = "Archiver.New"

	consumerConfig := sarama.NewConfig()

	client, err := sarama.NewConsumerGroup(conf.BrokerList, consumerGroup, consumerConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	consumer := &Consumer{
		ready:             make(chan bool),
		MessageAggregator: messageAggregator,
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
				log.Error().Err(err).Str("op", op).Msg("")
				return
			}

			if ctx.Err() != nil {
				return
			}
			a.GroupConsumer.ready = make(chan bool)
		}
	}()

	<-a.GroupConsumer.ready
}

func (a *Archiver) Close() {
	const op = "Archiver.Close"

	err := a.Client.Close()
	if err != nil {
		log.Error().Err(err).Str("op", op).Msg("")
	}
}
