package message

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

func (m *MessageUseCase) Read(ctx context.Context, userID uint64, message *entity.Message) error {
	if err := m.MessageStore.MarkRead(ctx, userID, message); err != nil {
		return err
	}

	return nil
}
