package messenger

import (
	"io"

	"github.com/kolesnikovm/messenger/proto"
)

func (s *Handler) SendMessage(stream proto.Messenger_SendMessageServer) error {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.Status{})
		}
		if err != nil {
			return err
		}

		m, err := convertPbToEntity(message)
		if err != nil {
			return err
		}

		err = s.Usecase.Send(stream.Context(), m)
		if err != nil {
			return err
		}
	}
}
