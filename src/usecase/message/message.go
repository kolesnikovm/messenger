package message

import "github.com/kolesnikovm/messenger/integration"

type MessageUseCase struct {
	messageSender integration.MessageSender
}

func New(messageSender integration.MessageSender) *MessageUseCase {
	return &MessageUseCase{
		messageSender: messageSender,
	}
}
