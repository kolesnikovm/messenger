package partitions

import (
	"github.com/golang/groupcache/consistenthash"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PartitionSet struct {
	hash         *consistenthash.Map
	partitionMap map[string]*pgxpool.Pool
}

const replicationFactor = 20

func New(partitions []*pgxpool.Pool) *PartitionSet {
	hash := consistenthash.New(len(partitions)*replicationFactor, nil)
	partitionMap := make(map[string]*pgxpool.Pool, len(partitions))

	for _, partition := range partitions {
		hash.Add(partition.Config().ConnString())
		partitionMap[partition.Config().ConnString()] = partition
	}

	return &PartitionSet{
		hash:         hash,
		partitionMap: partitionMap,
	}
}

func (p *PartitionSet) Close() {
	for _, partition := range p.partitionMap {
		partition.Close()
	}
}
