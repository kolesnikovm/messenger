package store

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

type Messages interface {
	BatchInsert(ctx context.Context, messages []*entity.Message) error
	GetMessageHistory(ctx context.Context, fromMessageID ulid.ULID, chatID string) ([]*entity.Message, error)
}
