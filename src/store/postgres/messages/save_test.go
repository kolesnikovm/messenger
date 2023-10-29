package messages

import (
	"context"
	"testing"
	"time"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := postgres.New(ctx, config.Postgres)
	require.NoError(t, err)

	messages := New(ctx, db, config.Postgres)

	for i := 0; i < 2*config.Postgres.BatchSize; i++ {
		messages.Save(&entity.Message{
			MessageID:   ulid.Make(),
			SenderID:    1,
			RecipientID: 2,
			Text:        "test",
		})
	}

	time.Sleep(5 * time.Second)
}
