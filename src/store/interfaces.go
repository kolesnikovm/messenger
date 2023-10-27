package store

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

type Messages interface {
	Save(ctx context.Context, message *entity.Message) error
}
