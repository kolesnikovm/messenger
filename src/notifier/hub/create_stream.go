package hub

import (
	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (s *StreamHub) CreateStream(recipientID string, sessionID ulid.ULID) <-chan *entity.Message {
	s.mx.Lock()
	defer s.mx.Unlock()

	var stream chan *entity.Message

	userStreams, ok := s.Streams[recipientID]
	if !ok {
		// TODO add coonfig for message buffer size
		stream = make(chan *entity.Message, 10)
		userStreams = map[ulid.ULID](chan *entity.Message){sessionID: stream}
		s.Streams[recipientID] = userStreams
	} else {
		userStreams[sessionID] = stream
	}

	return stream
}
