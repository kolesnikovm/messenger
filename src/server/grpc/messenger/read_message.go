package messenger

import (
	"context"

	"github.com/kolesnikovm/messenger/metrics"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *Handler) ReadMessage(ctx context.Context, req *proto.Message) (*proto.Message, error) {
	metrics.ActiveStreams.With(prometheus.Labels{"type": "read"}).Inc()
	defer metrics.ActiveStreams.With(prometheus.Labels{"type": "read"}).Dec()

	userID := ctx.Value(StringContextKey("userID")).(uint64)

	m, err := convertPbToEntity(req)
	if err != nil {
		return nil, err
	}

	err = s.Usecase.Read(ctx, userID, m)
	if err != nil {
		return nil, err
	}

	ack := &proto.Message{}

	return ack, nil
}
