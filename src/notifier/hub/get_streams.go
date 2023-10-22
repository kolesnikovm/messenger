package hub

import (
	"github.com/kolesnikovm/messenger/entity"
)

func (s *StreamHub) GetStreams(userID uint64) [](chan *entity.Message) {
	s.RLock()
	defer s.RUnlock()

	var userStreams [](chan *entity.Message)
	for _, stream := range s.Streams[userID] {
		userStreams = append(userStreams, stream)
	}

	return userStreams
}
