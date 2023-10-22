package kafka

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (k *KafkaMessageSender) Get(ctx context.Context, userID uint64, sessionID ulid.ULID) <-chan *entity.Message {
	var stream chan *entity.Message

	userStreams, ok := k.Streams[userID]
	if !ok {
		// TODO add coonfig for message buffer size
		stream = make(chan *entity.Message, 10)
		userStreams = map[ulid.ULID](chan *entity.Message){sessionID: stream}
		k.Streams[userID] = userStreams
	} else {
		userStreams[sessionID] = stream
	}

	go func() {
		<-ctx.Done()

		delete(userStreams, sessionID)
		// TODO use mutex
		if len(userStreams) == 0 {
			delete(k.Streams, userID)
		}
		close(stream)
	}()

	return stream
}
