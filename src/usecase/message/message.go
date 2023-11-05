package message

import (
	"github.com/kolesnikovm/messenger/notifier"
	"github.com/kolesnikovm/messenger/store"
)

type MessageUseCase struct {
	messageSender notifier.MessageSender
	messageStore  store.Messages
}

func New(messageSender notifier.MessageSender, messageStore store.Messages) *MessageUseCase {
	return &MessageUseCase{
		messageSender: messageSender,
		messageStore:  messageStore,
	}
}
