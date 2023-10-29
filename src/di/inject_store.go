package di

import (
	"context"

	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/kolesnikovm/messenger/store/postgres/messages"
)

func ProvideDB(ctx context.Context, conf configs.ServerConfig) (*postgres.DB, func(), error) {
	db, err := postgres.New(ctx, conf.Postgres)

	cleanup := func() {
		db.Close()
	}

	return db, cleanup, err
}

func ProvideMessages(ctx context.Context, db *postgres.DB, conf configs.ServerConfig) store.Messages {
	return messages.New(ctx, db, conf.Postgres)
}

var StoreSet = wire.NewSet(
	ProvideDB,
	ProvideMessages,
)
