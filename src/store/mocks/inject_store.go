package mocks

import (
	"testing"

	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/store"
)

func ProvideStore(t *testing.T) *MockMessages {
	return NewMockMessages(t)
}

var StoreSet = wire.NewSet(
	ProvideStore,
	wire.Bind(new(store.Messages), new(*MockMessages)),
)
