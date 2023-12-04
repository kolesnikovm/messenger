package migrations

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

const migrationsDir = "."

type Migrations struct {
	DB map[string]*sql.DB
}

func New(postgres *postgres.DB) (*Migrations, error) {
	const op = "migrations.New"

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db := map[string]*sql.DB{}
	for _, pool := range postgres.PartitionSet.GetAll() {
		key := fmt.Sprintf("%s:%d/%s",
			pool.Config().ConnConfig.Host,
			pool.Config().ConnConfig.Port,
			pool.Config().ConnConfig.Database,
		)
		db[key] = stdlib.OpenDBFromPool(pool)
	}

	return &Migrations{
		DB: db,
	}, nil
}

func (m *Migrations) Close() {
	for _, db := range m.DB {
		db.Close()
	}
}

func (m *Migrations) Run(command string, args ...string) error {
	const op = "migrations.Run"

	for url, db := range m.DB {
		log.Info().Msgf("run migration on %s", url)

		err := goose.Run(command, db, migrationsDir, args...)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
