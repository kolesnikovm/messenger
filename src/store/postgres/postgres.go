package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store/postgres/partitions"
)

type DB struct {
	PartitionSet *partitions.PartitionSet
}

func New(conf configs.Postgres) (*DB, error) {
	const op = "postgres.New"

	connectionPools := make([]*pgxpool.Pool, 0, len(conf.URL))

	for _, url := range conf.URL {
		pgxConf, err := pgxpool.ParseConfig(url)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if conf.MaxConns > 0 {
			pgxConf.MaxConns = conf.MaxConns
		}
		pgxConf.MaxConnLifetime = conf.MaxConnLifetime

		pool, err := pgxpool.NewWithConfig(context.Background(), pgxConf)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if err := pool.Ping(context.Background()); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		connectionPools = append(connectionPools, pool)
	}

	partitionSet := partitions.New(connectionPools)

	return &DB{
		PartitionSet: partitionSet,
	}, nil
}

func (d *DB) Close() {
	d.PartitionSet.Close()
}
