package migrations

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "."

type Migrations struct {
	DB *sql.DB
}

func New(postgres *postgres.DB) (*Migrations, error) {
	const op = "migrations.New"

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db := stdlib.OpenDBFromPool(postgres.Pool)

	return &Migrations{
		DB: db,
	}, nil
}

func (m *Migrations) Close() {
	m.DB.Close()
}

func (m *Migrations) Run(command string, args ...string) error {
	const op = "migrations.Run"

	err := goose.Run(command, m.DB, migrationsDir, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
