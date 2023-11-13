package message

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (m *MessageUseCase) GetHistory(ctx context.Context, chatID string, fromMessageID ulid.ULID, userID uint64, messageCount uint32, direction string) ([]*entity.Message, error) {
	const maxMessageCount = 1_000

	if messageCount > maxMessageCount {
		messageCount = maxMessageCount
	}

	messages, err := m.MessageStore.GetMessageHistory(ctx, fromMessageID, chatID, messageCount, direction)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
