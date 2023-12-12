package messenger

import (
	"io"

	"github.com/kolesnikovm/messenger/metrics"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *Handler) ReadMessage(stream proto.Messenger_ReadMessageServer) error {
	metrics.ActiveStreams.With(prometheus.Labels{"type": "read"}).Inc()
	defer metrics.ActiveStreams.With(prometheus.Labels{"type": "read"}).Dec()

	userID := stream.Context().Value(StringContextKey("userID")).(uint64)

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		m, err := convertPbToEntity(message)
		if err != nil {
			return err
		}

		msgID, err := s.Usecase.Read(stream.Context(), userID, m)
		if err != nil {
			return err
		}

		ack := &proto.Message{MessageID: msgID.String()}
		if err := stream.Send(ack); err != nil {
			return err
		}
	}
}
