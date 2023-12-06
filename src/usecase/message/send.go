package message

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (m *MessageUseCase) Send(ctx context.Context, message *entity.Message) (ulid.ULID, error) {
	msgID := message.MessageID
	chatID := message.GetChatID()

	lastOrderID, err := m.OrderIDCache.GetLastMessageOrderID(ctx, chatID)
	if err != nil {
		return msgID, err
	}

	if lastOrderID == 0 {
		lastOrderID, err = m.MessageStore.GetLastMessageOrderID(ctx, chatID)
		if err != nil {
			return msgID, err
		}

		err = m.OrderIDCache.SetLastMessageOrderID(ctx, chatID, lastOrderID)
		if err != nil {
			return msgID, err
		}
	}

	message.OrderID = lastOrderID + 1

	if err := m.MessageSender.Send(ctx, message); err != nil {
		return msgID, err
	}

	return msgID, nil
}
