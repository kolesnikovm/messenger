package tests

import (
	"context"
	"testing"
	"time"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/notifier/kafka"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestSendAndGet(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	kafkaMessageSender, err := kafka.New(&config.Kafka)
	require.NoError(t, err)
	defer kafkaMessageSender.Close()

	ctx := context.Background()
	kafkaMessageSender.Start(ctx)

	sessionID := ulid.Make()
	readCh, cleanup := kafkaMessageSender.Get(ctx, uint64(1), sessionID)
	defer cleanup()

	entityMessage := entity.NewMessage(ulid.Make(), 1, 2, "test")
	err = kafkaMessageSender.Send(ctx, entityMessage)
	require.NoError(t, err)

	for {
		select {
		case msg := <-readCh:
			require.Equal(t, entityMessage, msg)
			return
		case <-time.After(time.Second):
			t.Error("no message received")
			return
		}
	}
}
