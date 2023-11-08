package usecase

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

type Message interface {
	Send(ctx context.Context, message *entity.Message) error
	Get(ctx context.Context, userID uint64, sessionID ulid.ULID) (stream <-chan *entity.Message, cleanup func())
	GetHistory(ctx context.Context, chatID string, fromMessageID ulid.ULID, userID uint64) ([]*entity.Message, error)
}
