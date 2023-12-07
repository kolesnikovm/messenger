package messages

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kolesnikovm/messenger/entity"
)

const insert = "insert into messages (id, sender_id, chat_id, text) values ($1, $2, $3, $4) on conflict (id) do update set text = $4"

func (m *Messages) BatchInsert(ctx context.Context, messages []*entity.Message) error {
	const op = "Messages.BatchInsert"

	batches := map[*pgxpool.Pool]*pgx.Batch{}

	for _, msg := range messages {
		chatID := msg.GetChatID()
		partition := m.DB.PartitionSet.Get(chatID)

		batch, exists := batches[partition]
		if exists {
			batch.Queue(insert, msg.MessageID, msg.SenderID, chatID, msg.Text)
		} else {
			b := &pgx.Batch{}
			b.Queue(insert, msg.MessageID, msg.SenderID, chatID, msg.Text)

			batches[partition] = b
		}
	}

	for partition, batch := range batches {
		results := partition.SendBatch(ctx, batch)
		defer results.Close()

		for i := 0; i < batch.Len(); i++ {
			_, err := results.Exec()
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	return nil
}
