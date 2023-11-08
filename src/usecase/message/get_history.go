package message

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (m *MessageUseCase) GetHistory(ctx context.Context, chatID string, fromMessageID ulid.ULID, userID uint64) ([]*entity.Message, error) {
	messages, err := m.MessageStore.GetMessageHistory(ctx, fromMessageID, chatID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
