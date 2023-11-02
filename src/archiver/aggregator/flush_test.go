package aggregator

import (
	"testing"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/kolesnikovm/messenger/store/postgres/messages"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestFlush(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	db, err := postgres.New(config.Store.Postgres)
	require.NoError(t, err)

	messageStore := messages.New(db, config.Store.Postgres)

	messageAggregator := New(config.Store, messageStore)

	batchCount := 2
	flushCount := 0

	for i := 0; i < batchCount*config.Store.BatchSize; i++ {
		flushed := messageAggregator.Flush(&entity.Message{
			MessageID:   ulid.Make(),
			SenderID:    1,
			RecipientID: 2,
			Text:        "test",
		})

		if flushed {
			flushCount++
		}
	}

	require.Equal(t, batchCount, flushCount)
}
