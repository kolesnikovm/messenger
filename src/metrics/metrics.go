package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	MessagesReceivedTotal prometheus.Counter
	MessagesSentTotal     prometheus.Counter
	MessagesSavedTotal    prometheus.Counter
	HistoryRequestsTotal  prometheus.Counter
	ActiveStreams         *prometheus.GaugeVec
)

func init() {

	MessagesReceivedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_received_total",
			Help: "Total count of received messages",
		},
	)

	MessagesSentTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_sent_total",
			Help: "Total count of sent messages",
		},
	)

	MessagesSavedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_saved_total",
			Help: "Total count of saved messages",
		},
	)

	HistoryRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "history_requests_total",
			Help: "Total count of message history requests",
		},
	)

	ActiveStreams = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_streams",
			Help: "Number of active streams",
		},
		[]string{"type"},
	)

	prometheus.MustRegister(
		MessagesReceivedTotal,
		MessagesSentTotal,
		MessagesSavedTotal,
		ActiveStreams,
	)
}
