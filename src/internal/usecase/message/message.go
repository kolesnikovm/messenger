package message

import (
	"github.com/kolesnikovm/messenger/internal/entity"
	"github.com/rs/zerolog/log"
)

type Usecase struct{}

func New() *Usecase {
	return &Usecase{}
}

func (uc *Usecase) Send(message entity.Message) error {
	log.Info().Msgf("new message: %s", message.Text)
	return nil
}
