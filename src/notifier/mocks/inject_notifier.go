package mocks

import (
	"testing"

	"github.com/google/wire"
	"github.com/kolesnikovm/messenger/notifier"
)

func ProvideNotifier(t *testing.T) *MockMessageSender {
	return NewMockMessageSender(t)
}

var NotifierSet = wire.NewSet(
	ProvideNotifier,
	wire.Bind(new(notifier.MessageSender), new(*MockMessageSender)),
)
