package messages

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

type dbMessageID ulid.ULID

type dbMessage struct {
	MessageID dbMessageID
	SenderID  uint64
	ChatID    string
	Text      string
}

const (
	selectMessagesBackward = "select id, sender_id, chat_id, text from messages where chat_id = $1 and id < $2 order by id desc limit $3"
	selectMessagesForward  = "select id, sender_id, chat_id, text from messages where chat_id = $1 and id > $2 order by id asc  limit $3"
)

func (m *Messages) GetMessageHistory(ctx context.Context, fromMessageID ulid.ULID, chatID string, messageCount uint32, direction string) ([]*entity.Message, error) {
	const op = "Messages.GetMessageHistory"

	selectMessages := selectMessagesBackward
	if direction == "FORWARD" {
		selectMessages = selectMessagesForward
	}

	rows, err := m.DB.PartitionSet.Get(chatID).Query(ctx, selectMessages, chatID, fromMessageID, messageCount)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	messages := make([]*entity.Message, 0)
	for rows.Next() {
		m := &dbMessage{}

		if err := rows.Scan(&m.MessageID, &m.SenderID, &m.ChatID, &m.Text); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		message, err := m.getEntityMessage()
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (m *dbMessage) getEntityMessage() (*entity.Message, error) {
	const op = "getEntityMessage"

	recipientID, err := getRecipientID(m.ChatID, m.SenderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return entity.NewMessage(
		ulid.ULID(m.MessageID),
		m.SenderID,
		recipientID,
		m.Text,
	), nil
}

func getRecipientID(chatID string, senderID uint64) (uint64, error) {
	const op = "getRecipientID"

	switch entity.GetChatType(chatID) {
	case entity.Group:
		groupID, err := entity.GetGroupID(chatID)
		if err != nil {
			return 0, err
		}

		return groupID, nil
	case entity.Channel:
		channelID, err := entity.GetChannelID(chatID)
		if err != nil {
			return 0, err
		}

		return channelID, nil
	case entity.P2P:
		user1, user2, err := entity.GetUserIDs(chatID)
		if err != nil {
			return 0, err
		}

		if user1 == senderID {
			return user2, nil
		}

		return user1, nil
	}

	return 0, fmt.Errorf("%s: failed to determine recipient id", op)
}

func (id *dbMessageID) Scan(src interface{}) error {
	const op = "dbMessageID.Scan"

	switch x := src.(type) {
	case nil:
		return nil
	case string:
		msgID, err := uuid.Parse(x)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		*id = dbMessageID(msgID)

		return nil
	}

	return fmt.Errorf("%s: failed to determine interface type", op)
}
