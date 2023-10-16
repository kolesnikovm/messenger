package usecase

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

type Message interface {
	Send(context.Context, entity.Message) error
}
