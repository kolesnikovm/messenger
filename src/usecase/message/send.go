package message

import (
	"github.com/kolesnikovm/messenger/entity"
	"github.com/rs/zerolog/log"
)

func (uc *MessageUseCase) Send(message entity.Message) error {
	log.Info().Msgf("new message: %s", message.Text)
	return nil
}
