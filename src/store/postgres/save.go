package postgres

import (
	"context"
	"fmt"

	"github.com/kolesnikovm/messenger/entity"
)

func (m *Messages) Save(ctx context.Context, message *entity.Message) error {
	const op = "Messages.Save"

	insert := "insert into messages (id, sender_id, recipient_id, text) values ($1, $2, $3, $4)"

	_, err := m.db.Exec(ctx, insert, message.MessageID.String(), message.SenderID, message.RecipientID, message.Text)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
