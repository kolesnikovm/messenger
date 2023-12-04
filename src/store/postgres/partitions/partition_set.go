package partitions

import (
	"context"
	"fmt"

	"github.com/golang/groupcache/consistenthash"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PartitionSet struct {
	hash         *consistenthash.Map
	partitionMap map[string]*pgxpool.Pool
}

const replicationFactor = 20

func New(partitions []*pgxpool.Pool) (*PartitionSet, error) {
	hash := consistenthash.New(len(partitions)*replicationFactor, nil)
	partitionMap := make(map[string]*pgxpool.Pool, len(partitions))

	for _, partition := range partitions {
		dbID, err := getDatabaseID(partition)
		if err != nil {
			return nil, err
		}

		if p, exists := partitionMap[dbID]; exists {
			return nil, fmt.Errorf("cannot use same name for distinct partiotions %s and %s",
				getDBAddress(p),
				getDBAddress(partition))
		}
		partitionMap[dbID] = partition
		hash.Add(dbID)
	}

	return &PartitionSet{
		hash:         hash,
		partitionMap: partitionMap,
	}, nil
}

func (p *PartitionSet) Close() {
	for _, partition := range p.partitionMap {
		partition.Close()
	}
}

func getDatabaseID(db *pgxpool.Pool) (string, error) {
	var id string

	err := db.QueryRow(context.Background(), `select id from db_id`).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to get db id from %s: %w", getDBAddress(db), err)
	}

	return id, nil
}

func getDBAddress(db *pgxpool.Pool) string {
	return fmt.Sprintf("%s:%d/%s",
		db.Config().ConnConfig.Host,
		db.Config().ConnConfig.Port,
		db.Config().ConnConfig.Database)
}
