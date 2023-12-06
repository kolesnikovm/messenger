package messages

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const selectOrderID = "select order_id from messages where chat_id = $1 order by id desc limit 1"

func (m *Messages) GetLastMessageOrderID(ctx context.Context, chatID string) (uint64, error) {
	const op = "Messages.GetLastMessageOrderID"

	row := m.DB.PartitionSet.Get(chatID).QueryRow(ctx, selectOrderID, chatID)

	var orderID uint64

	if err := row.Scan(&orderID); err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return orderID, nil
}
