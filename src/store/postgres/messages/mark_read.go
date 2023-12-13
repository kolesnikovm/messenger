package messages

import (
	"context"
	"fmt"

	"github.com/kolesnikovm/messenger/entity"
)

const (
	markMessageRead = "update chat_participants set last_read_message = $1 where chat_id = $2 and user_id = $3"
)

func (m *Messages) MarkRead(ctx context.Context, userID uint64, message *entity.Message) error {
	const op = "Messages.MarkRead"

	chatID := message.GetChatID()

	if _, err := m.DB.PartitionSet.Get(chatID).Exec(ctx, markMessageRead, message.MessageID, chatID, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
