package messenger

import (
	"io"

	"github.com/kolesnikovm/messenger/metrics"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *Handler) SendMessage(stream proto.Messenger_SendMessageServer) error {
	metrics.ActiveStreams.With(prometheus.Labels{"type": "send"}).Inc()
	defer metrics.ActiveStreams.With(prometheus.Labels{"type": "send"}).Dec()

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

		msgID, err := s.Usecase.Send(stream.Context(), m)
		if err != nil {
			return err
		}

		ack := &proto.Message{MessageID: msgID.String()}
		if err := stream.Send(ack); err != nil {
			return err
		}

		metrics.MessagesReceivedTotal.Inc()
	}
}
