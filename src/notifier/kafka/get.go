package kafka

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (k *KafkaMessageSender) Get(ctx context.Context, recipientID string, sessionID ulid.ULID) (<-chan *entity.Message, func()) {
	stream := k.StreamHub.CreateStream(recipientID, sessionID)

	cleanup := func() {
		k.StreamHub.DeleteStream(recipientID, sessionID)
	}

	return stream, cleanup
}
