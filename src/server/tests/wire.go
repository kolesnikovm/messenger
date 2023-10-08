//go:build wireinject
// +build wireinject

package tests

import (
	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/di"
)

func InitializeSuite() *Suite {
	wire.Build(
		di.UsecaseSet,
		di.ServerSet,
		newSuite,
	)
	return &Suite{}
}
