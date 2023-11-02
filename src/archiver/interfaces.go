package archiver

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

type Archiver interface {
	Start(ctx context.Context)
}

type Aggregator interface {
	Flush(message *entity.Message) bool
}
