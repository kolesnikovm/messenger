package message

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (m *MessageUseCase) Send(ctx context.Context, message *entity.Message) (ulid.ULID, error) {
	msgID := message.MessageID

	if err := m.MessageSender.Send(ctx, message); err != nil {
		return msgID, err
	}

	return msgID, nil
}
