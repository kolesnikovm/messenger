package kafka

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

func (k *KafkaMessageSender) Get(ctx context.Context, userID uint64, deviceID int) <-chan *entity.Message {
	var stream chan *entity.Message

	userDeviceStreams, ok := k.UserStreams[userID]
	if !ok {
		// TODO add coonfig for message buffer size
		stream = make(chan *entity.Message, 10)
		userDeviceStreams = map[int](chan *entity.Message){deviceID: stream}
		k.UserStreams[userID] = userDeviceStreams
	} else {
		stream, ok := userDeviceStreams[deviceID]
		if !ok {
			stream = make(chan *entity.Message, 10)
		}
		userDeviceStreams[deviceID] = stream
	}

	go func() {
		<-ctx.Done()

		delete(userDeviceStreams, deviceID)
		// TODO use mutex
		if len(userDeviceStreams) == 0 {
			delete(k.UserStreams, userID)
		}
		close(stream)
	}()

	return stream
}
