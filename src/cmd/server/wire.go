//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/di"
)

func InitializeApplication(conf configs.ServerConfig) (*application, func(), error) {
	wire.Build(
		di.UsecaseSet,
		di.ServerSet,
		di.ProvideNotifier,
		newApplication,
	)
	return &application{}, nil, nil
}
