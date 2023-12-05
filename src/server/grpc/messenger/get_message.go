package messenger

import (
	"github.com/kolesnikovm/messenger/metrics"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/oklog/ulid/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func (h *Handler) GetMessage(msgRequest *proto.MessaggeRequest, stream proto.Messenger_GetMessageServer) error {
	userID := stream.Context().Value(StringContextKey("userID")).(uint64)

	sessionID := ulid.Make()

	messageCh, cleanup := h.Usecase.Get(stream.Context(), userID, sessionID)
	defer cleanup()

	metrics.ActiveStreams.With(prometheus.Labels{"type": "get"}).Inc()
	defer metrics.ActiveStreams.With(prometheus.Labels{"type": "get"}).Dec()

	for {
		select {
		case message := <-messageCh:
			protoMsg := convertEntityToPb(message)

			if err := stream.Send(protoMsg); err != nil {
				return err
			}

			metrics.MessagesSentTotal.Inc()
		case <-stream.Context().Done():
			return nil
		}
	}
}
