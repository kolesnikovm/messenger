package message

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

func (m *MessageUseCase) Send(ctx context.Context, message *entity.Message) error {
	if err := m.messageSender.Send(ctx, message); err != nil {
		return err
	}

	return nil
}
