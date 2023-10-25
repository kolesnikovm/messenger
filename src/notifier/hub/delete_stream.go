package hub

import (
	"github.com/oklog/ulid/v2"
)

func (s *StreamHub) DeleteStream(recipientID string, sessionID ulid.ULID) {
	s.mx.Lock()
	defer s.mx.Unlock()

	userStreams, ok := s.Streams[recipientID]
	if !ok {
		return
	}

	delete(userStreams, sessionID)
	if len(userStreams) == 0 {
		delete(s.Streams, recipientID)
	}
}
