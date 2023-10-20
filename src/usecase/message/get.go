package message

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

func (m *MessageUseCase) Get(ctx context.Context, userID uint64, deviceID int) <-chan *entity.Message {
	messageCh := m.messageSender.Get(ctx, userID, deviceID)

	return messageCh
}
