package hub

import (
	"github.com/kolesnikovm/messenger/entity"
)

func (s *StreamHub) GetStreams(recipientID string) [](chan *entity.Message) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	var userStreams [](chan *entity.Message)
	for _, stream := range s.Streams[recipientID] {
		userStreams = append(userStreams, stream)
	}

	return userStreams
}
