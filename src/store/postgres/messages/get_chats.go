package messages

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kolesnikovm/messenger/entity"
)

const (
	selectChats = "select chat_id from user_chats where user_id = $1"

	selectChatInfo = `select
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

	chatIDs := []string{}

	rows, err := m.DB.PartitionSet.Get(strconv.FormatUint(userID, 10)).Query(ctx, selectChats, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get user chats: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var chatID string

		if err := rows.Scan(&chatID); err != nil {
			return nil, fmt.Errorf("%s: failed to scan chat id: %w", op, err)
		}

		chatIDs = append(chatIDs, chatID)
	}
	rows.Close()

	chats := []*entity.Chat{}

	for _, chatID := range chatIDs {
		rows, err := m.DB.PartitionSet.Get(chatID).Query(ctx, selectChatInfo, userID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get chat info: %w", op, err)
		}
		defer rows.Close()

		for rows.Next() {
			chat := &entity.Chat{}

			if err := rows.Scan(&chat.ID, &chat.UnreadMessages); err != nil {
				return nil, fmt.Errorf("%s: failed to scan chat info: %w", op, err)
			}

			chats = append(chats, chat)
		}
	}

	return chats, nil
}
