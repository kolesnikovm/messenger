package message

import (
	"github.com/kolesnikovm/messenger/internal/entity"
)

type Message interface {
	Send(entity.Message) error
}
