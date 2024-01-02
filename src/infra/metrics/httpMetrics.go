package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"handler", "method", "code"},
	)
	httpResponseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_code",
			Help: "StatusCode of HTTP response",
		},
		[]string{"handler", "method", "code"},
	)
	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: []float64{0.01, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10, 20, 30, 40, 50, 60},
		}, []string{"handler", "method", "code"})

	httpConcurrentProcess = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_concurrent_process",
			Help: "Number http request being processed concurrently",
		})
)

func init() {
	prometheus.Register(httpRequests)
	prometheus.Register(httpResponseStatus)
	prometheus.Register(httpDuration)
	prometheus.Register(httpConcurrentProcess)
}
