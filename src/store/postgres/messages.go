package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kolesnikovm/messenger/configs"
)

type Messages struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, conf configs.Postgres) (*Messages, error) {
	const op = "postgres.New"

	pgxConf, err := pgxpool.ParseConfig(conf.URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if conf.MaxConns > 0 {
		pgxConf.MaxConns = conf.MaxConns
	}
	pgxConf.MaxConnLifetime = conf.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, pgxConf)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Messages{
		db: pool,
	}, nil
}

func (m *Messages) Close() {
	m.db.Close()
}
