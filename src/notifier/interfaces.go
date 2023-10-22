package notifier

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

type MessageSender interface {
	Send(context.Context, entity.Message) error
	Get(ctx context.Context, userID uint64, sessionID ulid.ULID) <-chan *entity.Message
}
