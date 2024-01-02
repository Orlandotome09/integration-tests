package metrics

import (
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	EventProcessStatusSuccess   = "success"
	EventProcessStatusError     = "error"
	EventProcessStatusPending   = "pending"
	EventProcessStatusDiscarded = "discarded"

	EventHandshakeStatusACK  = "ACK"
	EventHandshakeStatusNACK = "NACK"
)

type metricsGenerator struct {
	_interfaces.MetricsGenerator
}

func New() _interfaces.MetricsGenerator {

	return &metricsGenerator{}
}

func (ref *metricsGenerator) EventReceived(publishedTime time.Time) {

	eventWaitingDuration.Observe(time.Since(publishedTime).Seconds())
	eventConcurrentProcess.Inc()

}

func (ref *metricsGenerator) EventProcessed(handshake string, result string, requestReceived time.Time) {

	eventsProcessStatus.With(prometheus.Labels{"handshake": handshake, "status": result}).Inc()
	eventsProcessDuration.Observe(time.Since(requestReceived).Seconds())
	eventConcurrentProcess.Dec()

}
