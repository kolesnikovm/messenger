package messages

import (
	"github.com/kolesnikovm/messenger/store/postgres"
)

type Messages struct {
	DB *postgres.DB
}

func New(db *postgres.DB) *Messages {
	return &Messages{
		DB: db,
	}
}
