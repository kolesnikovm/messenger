package redis

import (
	"fmt"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
}

func New(conf configs.Redis) (*Cache, error) {
	const op = "redis.New"

	opts, err := redis.ParseURL(conf.URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rds := redis.NewClient(opts)

	return &Cache{
		Client: rds,
	}, nil
}

func (c *Cache) Close() {
	c.Client.Close()
}
