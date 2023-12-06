package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func (c *Cache) GetLastMessageOrderID(ctx context.Context, chatID string) (uint64, error) {
	const op = "Cache.GetLastMessageOrderID"

	val, err := c.Client.Get(ctx, chatID).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	orderID, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return orderID, nil
}
