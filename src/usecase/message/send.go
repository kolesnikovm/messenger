package message

import (
	"github.com/kolesnikovm/messenger/entity"
)

func (m *MessageUseCase) Send(message entity.Message) error {
	if err := m.messageSender.Send(message); err != nil {
		return err
	}

	return nil
}
