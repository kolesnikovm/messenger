package usecase

import (
	"github.com/kolesnikovm/messenger/entity"
)

type Message interface {
	Send(entity.Message) error
}
