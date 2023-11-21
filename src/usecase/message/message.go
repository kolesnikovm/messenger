package message

import (
	"github.com/kolesnikovm/messenger/notifier"
	"github.com/kolesnikovm/messenger/store"
)

type MessageUseCase struct {
	MessageSender notifier.MessageSender
	MessageStore  store.Messages
}
