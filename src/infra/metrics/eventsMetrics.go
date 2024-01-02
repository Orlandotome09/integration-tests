package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (

	// Pubsub metrics
	eventsProcessStatus = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "events_process_status",
		Help: "Number of events processed by status",
	}, []string{"handshake", "status"})

	eventsProcessDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "event_process_duration_seconds",
			Help:    "Duration of events processing",
			Buckets: []float64{0.01, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10, 20, 30, 40, 50, 60},
		})

	eventWaitingDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "event_waiting_duration_seconds",
		Help:    "Waiting time on queue before start processing",
		Buckets: []float64{0.01, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10, 20, 30, 40, 50, 60},
	})

	eventConcurrentProcess = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "event_concurrent_process",
		Help: "Number of events concurrently being processed",
	})
)

func init() {

	prometheus.Register(eventsProcessStatus)
	prometheus.Register(eventsProcessDuration)
	prometheus.Register(eventWaitingDuration)
	prometheus.Register(eventConcurrentProcess)
}
