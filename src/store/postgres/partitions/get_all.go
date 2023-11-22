package partitions

import "github.com/jackc/pgx/v5/pgxpool"

func (p *PartitionSet) GetAll() []*pgxpool.Pool {
	partitions := make([]*pgxpool.Pool, 0, len(p.partitionMap))

	for _, connectionPool := range p.partitionMap {
		partitions = append(partitions, connectionPool)
	}

	return partitions
}
