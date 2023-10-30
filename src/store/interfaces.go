package store

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

type Aggregator interface {
	Add(message *entity.Message)
}

type Messages interface {
	BatchInsert(ctx context.Context, messages []*entity.Message) error
}
