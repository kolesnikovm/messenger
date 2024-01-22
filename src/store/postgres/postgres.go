package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store/postgres/partitions"
	"github.com/rs/zerolog/log"
)

type DB struct {
	config      *configs.Postgres
	connections map[string]*pgxpool.Pool

	PartitionSet    *partitions.PartitionSet
	NewPartitionSet *partitions.PartitionSet
}

func New(conf *configs.Postgres) (*DB, error) {
	db := &DB{
		config:          conf,
		connections:     map[string]*pgxpool.Pool{},
		PartitionSet:    &partitions.PartitionSet{},
		NewPartitionSet: &partitions.PartitionSet{},
	}

	if err := db.createShards(); err != nil {
		return nil, err
	}

	return db, nil
}

func (d *DB) Close() {
	d.PartitionSet.Close()
	d.NewPartitionSet.Close()
}

func (d *DB) WatchResharding(ctx context.Context) {
	go func() {
		for {
			select {
			case <-d.config.Changed:
				if err := d.createShards(); err != nil {
					log.Error().Err(err).Msg("failed to recreate db shards")
					continue
				}

				if len(d.config.NewURL) == 0 {
					currentConnections := map[string]struct{}{}

					for _, url := range d.config.URL {
						currentConnections[url] = struct{}{}
					}

					for conn := range d.connections {
						if _, exists := currentConnections[conn]; !exists {
							delete(d.connections, conn)
						}
					}
				}

				log.Info().Msg("succesfully reconnected to databases")
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (d *DB) createShards() error {
	partitionSet, err := d.createConsistentHash(d.config.URL)
	if err != nil {
		return err
	}

	d.PartitionSet = partitionSet
	// atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&d.PartitionSet)), unsafe.Pointer(partitionSet))

	if len(d.config.NewURL) != 0 {
		newPartitionSet, err := d.createConsistentHash(d.config.NewURL)
		if err != nil {
			return err
		}

		d.NewPartitionSet = newPartitionSet
	}

	return nil
}

func (d *DB) createConsistentHash(URLs []string) (*partitions.PartitionSet, error) {
	connectionPools, err := d.connectDatabases(URLs)
	if err != nil {
		return nil, err
	}

	partitionSet, err := partitions.New(connectionPools)
	if err != nil {
		return nil, err
	}

	return partitionSet, nil
}

func (d *DB) connectDatabases(URLs []string) ([]*pgxpool.Pool, error) {
	connectionPools := make([]*pgxpool.Pool, 0, len(URLs))

	for _, url := range d.config.URL {
		pool, exists := d.connections[url]

		if !exists {
			var err error
			pool, err = createConnection(url, d.config)
			if err != nil {
				return nil, err
			}

			d.connections[url] = pool
		}

		connectionPools = append(connectionPools, pool)
	}

	return connectionPools, nil
}

func createConnection(url string, conf *configs.Postgres) (*pgxpool.Pool, error) {
	const op = "postgres.createConnection"

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

	return pool, nil
}
