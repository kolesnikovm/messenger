package messages

import (
	"context"
	"fmt"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

const selectMessages = "select id, sender_id, chat_id, text from messages where chat_id = $1 and id < $2 order by id desc limit $3"

func (m *Messages) GetMessageHistory(ctx context.Context, fromMessageID ulid.ULID, chatID string) ([]*entity.Message, error) {
	const op = "Messages.GetMessageHistory"

	// TODO add config
	messageCount := 50

	rows, err := m.DB.Query(ctx, selectMessages, chatID, fromMessageID, messageCount)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	messages := make([]*entity.Message, 0)
	for rows.Next() {
		m := &entity.Message{}

		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		recipientID, err := getRecipientID(values[2].(string), uint64(values[1].(int64)))
		if err != nil {
			return nil, err
		}
		m.RecipientID = recipientID

		messageIdBytes, ok := values[0].([16]byte)
		if !ok {
			return nil, fmt.Errorf("%s: failed to cast message id to bytes", op)
		}
		m.MessageID = ulid.ULID(messageIdBytes)

		senderID := uint64(values[1].(int64))
		m.SenderID = senderID

		text := values[3].(string)
		m.Text = text

		messages = append(messages, m)
	}

	return messages, nil
}

func getRecipientID(chatID string, senderID uint64) (uint64, error) {
	user1, user2, err := entity.ParseChatID(chatID)
	if err != nil {
		return 0, err
	}

	if user1 == senderID {
		return user2, nil
	}

	return user1, nil
}
