package hub

import (
	"github.com/kolesnikovm/messenger/entity"
)

func (s *StreamHub) GetStreams(recepientID string) [](chan *entity.Message) {
	s.RLock()
	defer s.RUnlock()

	var userStreams [](chan *entity.Message)
	for _, stream := range s.Streams[recepientID] {
		userStreams = append(userStreams, stream)
	}

	return userStreams
}
