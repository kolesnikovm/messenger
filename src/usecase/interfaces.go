package usecase

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

type Message interface {
	Send(context.Context, entity.Message) error
	Get(ctx context.Context, userID string, sessionID ulid.ULID, chatID string) (stream <-chan *entity.Message, cleanup func(), err error)
}
