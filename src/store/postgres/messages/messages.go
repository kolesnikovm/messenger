package messages

import (
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store/postgres"
)

type Messages struct {
	DB     *postgres.DB
	Config configs.Postgres
}

func New(db *postgres.DB, conf configs.Postgres) *Messages {
	return &Messages{
		DB:     db,
		Config: conf,
	}
}
