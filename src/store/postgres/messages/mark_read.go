package messages

import (
	"context"
	"fmt"

	"github.com/kolesnikovm/messenger/entity"
)

const (
	markMessageRead = `insert into chat_participants (user_id, chat_id, last_read_message) values ($1, $2, $3)
	on conflict (user_id, chat_id) do update set last_read_message = $3`
)

func (m *Messages) MarkRead(ctx context.Context, userID uint64, message *entity.Message) error {
	const op = "Messages.MarkRead"

	chatID := message.GetChatID()

	if _, err := m.DB.PartitionSet.Get(chatID).Exec(ctx, markMessageRead, userID, chatID, message.MessageID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
