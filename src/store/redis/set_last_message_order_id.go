package redis

import (
	"context"
	"fmt"
	"strconv"
)

func (c *Cache) SetLastMessageOrderID(ctx context.Context, chatID string, orderID uint64) error {
	const op = "Cache.SetLastMessageOrderID"

	val := strconv.FormatUint(orderID, 10)

	err := c.Client.Set(ctx, chatID, val, 0).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
