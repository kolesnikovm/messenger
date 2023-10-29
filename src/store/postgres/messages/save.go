package messages

import (
	"github.com/kolesnikovm/messenger/entity"
)

func (m *Messages) Save(message *entity.Message) {
	m.Messages <- message
}
