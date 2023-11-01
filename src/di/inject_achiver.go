package di

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/archiver"
	"github.com/kolesnikovm/messenger/archiver/kafka"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store"
	"github.com/kolesnikovm/messenger/store/aggregator"
)

func ProvideAggregator(conf configs.ServerConfig, messageStore store.Messages) store.Aggregator {
	return aggregator.New(conf.Store, messageStore)
}

func ProvideArchiver(conf configs.ServerConfig, messageAggregator store.Aggregator) (archiver.Archiver, func(), error) {
	archiver, err := kafka.New(conf.KafkaConfig, messageAggregator)

	cleanup := func() {
		archiver.Close()
	}

	return archiver, cleanup, err
}

var ArchiverSet = wire.NewSet(
	ProvideAggregator,
	ProvideArchiver,
)
