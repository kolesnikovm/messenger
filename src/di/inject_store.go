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

var StoreSet = wire.NewSet(
	ProvideDB,
	messages.New,
	wire.Bind(new(store.Messages), new(*messages.Messages)),
)
