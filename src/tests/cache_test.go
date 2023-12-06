package tests

import (
	"context"
	"testing"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/store/redis"
	"github.com/stretchr/testify/require"
)

func TestGetOrderID(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	cache, err := redis.New(config.Redis)
	require.NoError(t, err)

	ctx := context.Background()
	orderID := uint64(1)

	err = cache.SetLastMessageOrderID(ctx, "1:2", orderID)
	require.NoError(t, err)

	cacheOrderID, err := cache.GetLastMessageOrderID(ctx, "1:2")
	require.NoError(t, err)

	require.Equal(t, orderID, cacheOrderID)
}
