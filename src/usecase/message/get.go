package message

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (m *MessageUseCase) Get(ctx context.Context, userID uint64, sessionID ulid.ULID) (<-chan *entity.Message, func()) {
	messageCh, cleanup := m.MessageSender.Get(ctx, userID, sessionID)

	return messageCh, cleanup
}
