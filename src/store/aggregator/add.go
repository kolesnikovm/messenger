package aggregator

import (
	"github.com/kolesnikovm/messenger/entity"
)

func (a *Aggregator) Add(message *entity.Message) {
	a.Messages <- message
}
