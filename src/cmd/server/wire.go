//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/di"
)

func InitializeApplication(conf configs.ServerConfig) *application {
	wire.Build(
		di.UsecaseSet,
		di.ServerSet,
		di.IntegrationSet,
		newApplication,
	)
	return &application{}
}
