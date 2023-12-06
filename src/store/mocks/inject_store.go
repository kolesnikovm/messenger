package mocks

import (
	"testing"

	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/store"
)

func ProvideStore(t *testing.T) *MockMessages {
	return NewMockMessages(t)
}

func ProvideCache(t *testing.T) *MockOrderIDCacher {
	return NewMockOrderIDCacher(t)
}

var StoreSet = wire.NewSet(
	ProvideStore,
	wire.Bind(new(store.Messages), new(*MockMessages)),
	ProvideCache,
	wire.Bind(new(store.OrderIDCacher), new(*MockOrderIDCacher)),
)
