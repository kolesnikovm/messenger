package messenger

import (
	"io"

	"github.com/kolesnikovm/messenger/proto"
	"github.com/rs/zerolog/log"
)

func (s *Handler) SendMessage(stream proto.Messenger_SendMessageServer) error {
	errorCount := int32(0)
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.Status{
				ErrorCount: errorCount,
			})
		}
		if err != nil {
			log.Error().Err(err).Msg("")
			return err
		}

		m := s.transformMessageRPC(message)

		err = s.Usecase.Send(m)
		if err != nil {
			log.Error().Err(err).Msg("failed to send message")
			errorCount++
		}
	}
}
