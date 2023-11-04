package messages

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kolesnikovm/messenger/entity"
)

const insert = "insert into messenger.messages (id, sender_id, chat_id, text) values ($1, $2, $3, $4) on conflict (id) do update set text = $4"

func (m *Messages) BatchInsert(ctx context.Context, messages []*entity.Message) error {
	const op = "Messages.BatchInsert"

	batch := &pgx.Batch{}

	for _, msg := range messages {
		chatID := msg.GetChatID()
		batch.Queue(insert, msg.MessageID, msg.SenderID, chatID, msg.Text)
	}

	results := m.DB.SendBatch(ctx, batch)

	for i := 0; i < batch.Len(); i++ {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return results.Close()
}
