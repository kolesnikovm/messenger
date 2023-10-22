package hub

import (
	"sync"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

type StreamHub struct {
	sync.RWMutex
	Streams map[uint64]map[ulid.ULID](chan *entity.Message)
}

func New() *StreamHub {
	return &StreamHub{
		Streams: make(map[uint64]map[ulid.ULID](chan *entity.Message)),
	}
}
