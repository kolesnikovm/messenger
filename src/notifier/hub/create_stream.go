package hub

import (
	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (s *StreamHub) CreateStream(userID uint64, sessionID ulid.ULID) <-chan *entity.Message {
	s.Lock()
	defer s.Unlock()

	var stream chan *entity.Message

	userStreams, ok := s.Streams[userID]
	if !ok {
		// TODO add coonfig for message buffer size
		stream = make(chan *entity.Message, 10)
		userStreams = map[ulid.ULID](chan *entity.Message){sessionID: stream}
		s.Streams[userID] = userStreams
	} else {
		userStreams[sessionID] = stream
	}

	return stream
}
