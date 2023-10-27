//go:build wireinject
// +build wireinject

package server

import (
	"context"

	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/di"
)

func InitializeApplication(ctx context.Context, conf configs.ServerConfig) (*application, func(), error) {
	wire.Build(
		di.UsecaseSet,
		di.ServerSet,
		di.ProvideNotifier,
		di.ProvideStore,
		newApplication,
	)
	return &application{}, nil, nil
}
