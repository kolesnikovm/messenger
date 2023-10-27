package di

import (
	"context"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store"
	"github.com/kolesnikovm/messenger/store/postgres"
)

func ProvideStore(ctx context.Context, conf configs.ServerConfig) (store.Messages, func(), error) {
	messageStore, err := postgres.New(ctx, conf.Postgres)

	cleanup := func() {
		messageStore.Close()
	}

	return messageStore, cleanup, err
}
