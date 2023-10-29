package store

import (
	"github.com/kolesnikovm/messenger/entity"
)

type Messages interface {
	Save(message *entity.Message)
}
