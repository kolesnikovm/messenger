package di

import (
	"context"

	"github.com/kolesnikovm/messenger/archiver"
	"github.com/kolesnikovm/messenger/archiver/kafka"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store"
)

func ProvideArchiver(ctx context.Context, conf configs.ServerConfig, messageAggregator store.Aggregator) (archiver.Archiver, func(), error) {
	archiver, err := kafka.New(ctx, conf.KafkaConfig, messageAggregator)

	cleanup := func() {
		archiver.Close()
	}

	return archiver, cleanup, err
}
