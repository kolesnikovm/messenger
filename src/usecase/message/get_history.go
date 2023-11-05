package message

import (
	"context"
	"fmt"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (m *MessageUseCase) GetHistory(ctx context.Context, chatID string, fromMessageID ulid.ULID, userID uint64) ([]*entity.Message, error) {
	const op = "MessageUseCase.GetGistory"

	user1, user2, err := entity.ParseChatID(chatID)
	if err != nil {
		return nil, err
	}

	if user1 != userID && user2 != userID {
		return nil, fmt.Errorf("%s: permission denied for user %d on chat %s", op, userID, chatID)
	}

	messages, err := m.messageStore.GetMessageHistory(ctx, fromMessageID, chatID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
