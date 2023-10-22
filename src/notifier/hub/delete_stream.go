package hub

import (
	"github.com/oklog/ulid/v2"
)

func (s *StreamHub) DeleteStream(userID uint64, sessionID ulid.ULID) {
	s.Lock()
	defer s.Unlock()

	userStreams, ok := s.Streams[userID]
	if !ok {
		return
	}

	delete(userStreams, sessionID)
	if len(userStreams) == 0 {
		delete(s.Streams, userID)
	}
}
