package kafka

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (k *KafkaMessageSender) Get(ctx context.Context, userID uint64, sessionID ulid.ULID) <-chan *entity.Message {
	stream := k.StreamHub.CreateStream(userID, sessionID)

	go func() {
		<-ctx.Done()

		k.StreamHub.DeleteStream(userID, sessionID)
	}()

	return stream
}
