package store

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

type Messages interface {
	BatchInsert(ctx context.Context, messages []*entity.Message) error
}
