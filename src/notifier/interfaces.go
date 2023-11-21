package notifier

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

type MessageSender interface {
	Send(ctx context.Context, message *entity.Message) error
	Get(ctx context.Context, userID uint64, sessionID ulid.ULID) (stream <-chan *entity.Message, cleanup func())
	Start(ctx context.Context)
}
