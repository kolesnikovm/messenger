package notifier

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

type MessageSender interface {
	Send(context.Context, entity.Message) error
	Get(ctx context.Context, recepientID string, sessionID ulid.ULID) (stream <-chan *entity.Message, cleanup func())
}
