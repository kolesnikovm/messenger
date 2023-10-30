package aggregator

import (
	"context"
	"testing"
	"time"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/kolesnikovm/messenger/store/postgres/messages"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := postgres.New(config.Store.Postgres)
	require.NoError(t, err)

	messageStore := messages.New(db, config.Store.Postgres)

	messageAggregator := New(config.Store, messageStore)
	messageAggregator.Start(ctx)

	for i := 0; i < 2*config.Store.BatchSize; i++ {
		messageAggregator.Add(&entity.Message{
			MessageID:   ulid.Make(),
			SenderID:    1,
			RecipientID: 2,
			Text:        "test",
		})
	}

	time.Sleep(5 * time.Second)
}
