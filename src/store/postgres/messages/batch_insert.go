package messages

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/rs/zerolog/log"
)

const (
	insertMessage = "insert into messages (id, sender_id, chat_id, text, order_id) values ($1, $2, $3, $4, $5) on conflict (id) do update set text = $4"

	selectMessageCounter = "select message_counter from chats where id = $1"
	updateMessageCounter = "insert into chats (id, message_counter) values ($1, $2) on conflict (id) do update set message_counter = $2"
)

func (m *Messages) BatchInsert(ctx context.Context, messages []*entity.Message) error {
	const op = "Messages.BatchInsert"

	batches := map[*pgxpool.Pool]*pgx.Batch{}

	transactions := map[*pgxpool.Pool]pgx.Tx{}
	defer func() {
		for _, tx := range transactions {
			if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
				log.Err(err).Msgf("%s: failed to rollback transaction", op)
			}
		}
	}()

	for _, msg := range messages {
		chatID := msg.GetChatID()
		partition := m.DB.PartitionSet.Get(chatID)

		tx, exists := transactions[partition]
		if !exists {
			var err error
			tx, err = partition.Begin(ctx)
			if err != nil {
				return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
			}

			transactions[partition] = tx
		}

		orderID, err := m.getMessageCounter(ctx, tx, chatID)
		if err != nil {
			return err
		}
		orderID++

		err = m.updateMessageCounter(ctx, tx, chatID, orderID)
		if err != nil {
			return err
		}

		batch, exists := batches[partition]
		if exists {
			batch.Queue(insertMessage, msg.MessageID, msg.SenderID, chatID, msg.Text, orderID)
		} else {
			b := &pgx.Batch{}
			b.Queue(insertMessage, msg.MessageID, msg.SenderID, chatID, msg.Text, orderID)

			batches[partition] = b
		}
	}

	for partition, tx := range transactions {
		batch := batches[partition]

		results := tx.SendBatch(ctx, batch)
		defer results.Close()

		for i := 0; i < batch.Len(); i++ {
			_, err := results.Exec()
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

		results.Close()

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
		}
	}

	return nil
}

func (m *Messages) getMessageCounter(ctx context.Context, tx pgx.Tx, chatID string) (uint64, error) {
	const op = "Messages.getMessageCounter"

	m.chats.mu.RLock()
	orderID, exists := m.chats.counters[chatID]
	m.chats.mu.RUnlock()

	if !exists {
		row := tx.QueryRow(ctx, selectMessageCounter, chatID)

		if err := row.Scan(&orderID); err != nil {
			if err == pgx.ErrNoRows {
				return 0, nil
			}
			return 0, fmt.Errorf("%s: %w", op, err)
		}

	}

	return orderID, nil
}

func (m *Messages) updateMessageCounter(ctx context.Context, tx pgx.Tx, chatID string, orderID uint64) error {
	const op = "Messages.updateMessageCounter"

	m.chats.mu.Lock()
	m.chats.counters[chatID] = orderID
	m.chats.mu.Unlock()

	_, err := tx.Exec(ctx, updateMessageCounter, chatID, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
