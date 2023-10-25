package kafka

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (k *KafkaMessageSender) Get(ctx context.Context, userID uint64, sessionID ulid.ULID) (<-chan *entity.Message, func()) {
	stream := k.StreamHub.CreateStream(userID, sessionID)

	cleanup := func() {
		k.StreamHub.DeleteStream(userID, sessionID)
	}

	return stream, cleanup
}
