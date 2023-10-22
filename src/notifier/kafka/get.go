package kafka

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (k *KafkaMessageSender) Get(ctx context.Context, recepientID string, sessionID ulid.ULID) (<-chan *entity.Message, func()) {
	stream := k.StreamHub.CreateStream(recepientID, sessionID)

	cleanup := func() {
		k.StreamHub.DeleteStream(recepientID, sessionID)
	}

	return stream, cleanup
}
