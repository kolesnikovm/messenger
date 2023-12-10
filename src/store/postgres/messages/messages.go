package messages

import (
	"sync"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store/postgres"
)

type Messages struct {
	DB     *postgres.DB
	Config configs.Postgres
	chats  *messageCounter
}

type messageCounter struct {
	mu       sync.RWMutex
	counters map[string]uint64
}

func New(db *postgres.DB, conf configs.Postgres) *Messages {
	return &Messages{
		DB:     db,
		Config: conf,
		chats: &messageCounter{
			counters: map[string]uint64{},
		},
	}
}
