//go:build wireinject
// +build wireinject

package tests

import (
	"testing"

	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/di"
	notifier "github.com/kolesnikovm/messenger/notifier/mocks"
	store "github.com/kolesnikovm/messenger/store/mocks"
)

func InitializeSuite(t *testing.T, conf configs.ServerConfig) (*Suite, error) {
	wire.Build(
		di.UsecaseSet,
		di.ServerSet,
		notifier.NotifierSet,
		store.StoreSet,
		newConnection,
		newClient,
		wire.Struct(new(Suite), "*"),
	)
	return &Suite{}, nil
}
