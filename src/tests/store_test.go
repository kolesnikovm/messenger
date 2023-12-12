package tests

import (
	"context"
	"testing"
	"time"

	"github.com/kolesnikovm/messenger/archiver/kafka"
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

	messageArchiver, err := kafka.New(config.Kafka, config.Archiver, messageStore)
	require.NoError(t, err)

	message1 := entity.NewMessage(ulid.Make(), 1, 2, "test")

	time.Sleep(1 * time.Second)
	message2 := entity.NewMessage(ulid.Make(), 1, 2, "test")
	messages := []*entity.Message{message1, message2}

	ctx := context.Background()
	messageArchiver.GroupConsumer.SendMessages(ctx, messages)

	time.Sleep(1 * time.Second)

	historyMessages, err := messageStore.GetMessageHistory(ctx, message2.MessageID, "1:2", 50, "BACKWARD")
	require.NoError(t, err)

	require.Contains(t, historyMessages, message1)
}

func TestGetHistoryForward(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	db, err := postgres.New(config.Postgres)
	require.NoError(t, err)

	messageStore := messages.New(db, config.Postgres)

	messageArchiver, err := kafka.New(config.Kafka, config.Archiver, messageStore)
	require.NoError(t, err)

	message1 := entity.NewMessage(ulid.Make(), 1, 2, "test")
	time.Sleep(1 * time.Second)
	message2 := entity.NewMessage(ulid.Make(), 1, 2, "test")
	messages := []*entity.Message{message1, message2}

	ctx := context.Background()
	messageArchiver.GroupConsumer.SendMessages(ctx, messages)

	time.Sleep(1 * time.Second)

	historyMessages, err := messageStore.GetMessageHistory(ctx, message1.MessageID, "1:2", 50, "FORWARD")
	require.NoError(t, err)

	require.Contains(t, historyMessages, message2)
}

func TestMarkRead(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	db, err := postgres.New(config.Postgres)
	require.NoError(t, err)

	messageStore := messages.New(db, config.Postgres)

	ctx := context.Background()
	message := entity.NewMessage(ulid.Make(), 1, 2, "test")

	err = messageStore.MarkRead(ctx, 1, message)
	require.NoError(t, err)
}
