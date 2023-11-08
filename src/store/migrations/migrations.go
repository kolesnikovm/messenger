package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var embedMigrations embed.FS

type Migrations struct {
	DB *sql.DB
}

func New(postgres *postgres.DB) (*Migrations, error) {
	const op = "migrations.New"

	goose.SetBaseFS(embedMigrations)

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
