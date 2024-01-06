package message

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

func (m *MessageUseCase) GetChats(ctx context.Context, userID uint64) ([]*entity.Chat, error) {
	chats, err := m.MessageStore.GetChats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return chats, nil
}
