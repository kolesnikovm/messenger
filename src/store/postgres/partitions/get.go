package partitions

import "github.com/jackc/pgx/v5/pgxpool"

func (p *PartitionSet) Get(key string) *pgxpool.Pool {
	partition := p.hash.Get(key)

	return p.partitionMap[partition]
}
