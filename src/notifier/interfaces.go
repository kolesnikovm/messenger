package notifier

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

type MessageSender interface {
	Send(context.Context, entity.Message) error
}
