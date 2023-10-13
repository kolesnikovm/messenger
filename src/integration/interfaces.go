package integration

import "github.com/kolesnikovm/messenger/entity"

type MessageSender interface {
	Send(entity.Message) error
}
