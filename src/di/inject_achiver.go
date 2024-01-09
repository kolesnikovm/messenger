package di

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/archiver"
	"github.com/kolesnikovm/messenger/archiver/kafka"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store"
)

func ProvideArchiver(conf *configs.ServerConfig, messageStore store.Messages) (archiver.Archiver, func(), error) {
	archiver, err := kafka.New(&conf.Kafka, &conf.Archiver, messageStore)

	cleanup := func() {
		archiver.Close()
	}

	return archiver, cleanup, err
}

var ArchiverSet = wire.NewSet(
	ProvideArchiver,
)
