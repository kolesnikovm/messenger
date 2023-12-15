package messages

import (
	"context"
	"fmt"

	"github.com/kolesnikovm/messenger/entity"
)

const (
	selectChats = `select
		c.id,
		case
			when last_read_message is null then c.message_counter
			when last_read_message is not null then c.message_counter - m.order_id
		end as unread_messages
	from
		chats c
	join chat_participants cp on
		c.id = cp.chat_id
	left join messages m on
		m.id = cp.last_read_message
	where
		cp.user_id = $1`
)

func (m *Messages) GetChats(ctx context.Context, userID uint64) ([]*entity.Chat, error) {
	const op = "Messages.GetChats"

	chats := []*entity.Chat{}

	for _, partition := range m.DB.PartitionSet.GetAll() {
		rows, err := partition.Query(ctx, selectChats, userID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()

		for rows.Next() {
			chat := &entity.Chat{}

			if err := rows.Scan(&chat.ID, &chat.UnreadMessages); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}

			chats = append(chats, chat)
		}
	}

	return chats, nil
}
