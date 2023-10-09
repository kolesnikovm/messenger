//go:build wireinject
// +build wireinject

package tests

import (
	"testing"

	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/di"
)

func InitializeSuite(t *testing.T) (*Suite, error) {
	wire.Build(
		di.UsecaseSet,
		di.ServerSet,
		newSuite,
	)
	return &Suite{}, nil
}
