package message

import "github.com/kolesnikovm/messenger/notifier"

type MessageUseCase struct {
	messageSender notifier.MessageSender
}

func New(messageSender notifier.MessageSender) *MessageUseCase {
	return &MessageUseCase{
		messageSender: messageSender,
	}
}
