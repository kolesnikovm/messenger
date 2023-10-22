package hub

import (
	"github.com/kolesnikovm/messenger/entity"
)

func (s *StreamHub) GetStreams(recipientID string) [](chan *entity.Message) {
	s.RLock()
	defer s.RUnlock()

	var userStreams [](chan *entity.Message)
	for _, stream := range s.Streams[recipientID] {
		userStreams = append(userStreams, stream)
	}

	return userStreams
}
