package message

import (
	"context"

	"github.com/kolesnikovm/messenger/entity"
)

func (m *MessageUseCase) Send(message entity.Message) error {
	ctx := context.Background()
	if err := m.messageSender.Send(ctx, message); err != nil {
		return err
	}

	return nil
}
