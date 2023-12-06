package di

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/kolesnikovm/messenger/store/postgres/messages"
	"github.com/kolesnikovm/messenger/store/redis"
)

func ProvideDB(conf configs.ServerConfig) (*postgres.DB, func(), error) {
	db, err := postgres.New(conf.Postgres)

	cleanup := func() {
		db.Close()
	}

	return db, cleanup, err
}

func ProvideMessages(db *postgres.DB, conf configs.ServerConfig) store.Messages {
	return messages.New(db, conf.Postgres)
}

func ProvideCache(conf configs.ServerConfig) (store.OrderIDCacher, func(), error) {
	cache, err := redis.New(conf.Redis)

	cleanup := func() {
		cache.Close()
	}

	return cache, cleanup, err
}

var StoreSet = wire.NewSet(
	ProvideDB,
	ProvideMessages,
	ProvideCache,
)
