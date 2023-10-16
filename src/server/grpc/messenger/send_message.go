package messenger

import (
	"io"

	"github.com/kolesnikovm/messenger/proto"
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
			return err
		}

		m := s.transformMessageRPC(message)

		err = s.Usecase.Send(stream.Context(), m)
		if err != nil {
			errorCount++
			return err
		}
	}
}
