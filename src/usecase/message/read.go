package message

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (m *MessageUseCase) Read(ctx context.Context, userID uint64, message *entity.Message) (ulid.ULID, error) {
	msgID := message.MessageID

	if err := m.MessageStore.MarkRead(ctx, userID, message); err != nil {
		return msgID, err
	}

	return msgID, nil
}
