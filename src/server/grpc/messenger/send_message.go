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

		metrics.MessagesReceivedTotal.Inc()
	}
}
