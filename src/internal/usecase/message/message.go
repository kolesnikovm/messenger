package message

import (
	"github.com/kolesnikovm/messenger/internal/entity"
	"github.com/rs/zerolog/log"
)

type MessageUseCase struct{}

func New() *MessageUseCase {
	return &MessageUseCase{}
}

func (uc *MessageUseCase) Send(message entity.Message) error {
	log.Info().Msgf("new message: %s", message.Text)
	return nil
}
