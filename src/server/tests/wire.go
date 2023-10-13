//go:build wireinject
// +build wireinject

package tests

import (
	"testing"

	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/di"
)

func InitializeSuite(t *testing.T, conf configs.ServerConfig) (*Suite, error) {
	wire.Build(
		di.UsecaseSet,
		di.ServerSet,
		di.IntegrationSet, // TODO replace with mock
		newSuite,
	)
	return &Suite{}, nil
}
