package di

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store"
	"github.com/kolesnikovm/messenger/store/aggregator"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/kolesnikovm/messenger/store/postgres/messages"
)

func ProvideDB(conf configs.ServerConfig) (*postgres.DB, func(), error) {
	db, err := postgres.New(conf.Store.Postgres)

	cleanup := func() {
		db.Close()
	}

	return db, cleanup, err
}

func ProvideMessages(db *postgres.DB, conf configs.ServerConfig) store.Messages {
	return messages.New(db, conf.Store.Postgres)
}

func ProvideAggregator(conf configs.ServerConfig, messageStore store.Messages) store.Aggregator {
	return aggregator.New(conf.Store, messageStore)
}

var StoreSet = wire.NewSet(
	ProvideDB,
	ProvideMessages,
	ProvideAggregator,
)
