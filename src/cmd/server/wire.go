//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/di"
)

func InitializeApplication() *application {
	wire.Build(
		di.UsecaseSet,
		di.ServerSet,
		newApplication,
	)
	return &application{}
}
