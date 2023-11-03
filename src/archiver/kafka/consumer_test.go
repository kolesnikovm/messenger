package kafka

import (
	"context"
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

	db, err := postgres.New(config.Postgres)
	require.NoError(t, err)

	messageStore := messages.New(db, config.Postgres)

	messageArchiver, err := New(config.KafkaConfig, config.Archiver, messageStore)
	require.NoError(t, err)

	messages := []*entity.Message{{
		MessageID:   ulid.Make(),
		SenderID:    1,
		RecipientID: 2,
		Text:        "test",
	}}

	messageArchiver.GroupConsumer.sendMessages(context.Background(), messages)
}
