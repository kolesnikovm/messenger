package kafka

import (
	"context"
	"testing"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	kafkaMessageSender, err := New(config.KafkaConfig)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	sessionID := ulid.Make()
	readCh := kafkaMessageSender.Get(ctx, uint64(1), sessionID)

	entityMessage := &entity.Message{
		MessageID:   ulid.Make(),
		SenderID:    1,
		RecipientID: 2,
		Text:        "test",
	}

	go func() {
		userStreams := kafkaMessageSender.StreamHub.GetStreams(uint64(1))
		require.Equal(t, 1, len(userStreams))

		for _, stream := range userStreams {
			for i := 0; i < 10; i++ {
				stream <- entityMessage
			}
		}
	}()

	go func() {
		cancel()
		kafkaMessageSender.Close()
	}()

	for {
		select {
		case msg := <-readCh:
			require.Equal(t, entityMessage, msg)
		case <-ctx.Done():
			return
		}
	}
}
