package message

import (
	"time"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/infra/metrics"
	"cloud.google.com/go/pubsub"
)

type Message interface {
	ID() string
	Data() []byte
	PublishTime() time.Time
	Ack()
	Nack(result string) bool
	ConcurrencyKey() string
}

type pubsubWrapper struct {
	m       *pubsub.Message
	metrics interfaces.MetricsGenerator
	start   time.Time
	ack     bool
}

func New(m *pubsub.Message, metrics interfaces.MetricsGenerator) Message {
	metrics.EventReceived(m.PublishTime)
	return &pubsubWrapper{
		m:       m,
		metrics: metrics,
		start:   time.Now(),
	}
}

func (ref *pubsubWrapper) ID() string {
	return ref.m.ID
}

func (ref *pubsubWrapper) Data() []byte {
	return ref.m.Data
}

func (ref *pubsubWrapper) PublishTime() time.Time {
	return ref.m.PublishTime
}

func (ref *pubsubWrapper) Ack() {
	ref.metrics.EventProcessed(metrics.EventHandshakeStatusACK, metrics.EventProcessStatusSuccess, ref.start)
	ref.m.Ack()
	ref.ack = true
}

func (ref *pubsubWrapper) Nack(result string) bool {
	if ref.ack {
		return false
	}
	ref.metrics.EventProcessed(metrics.EventHandshakeStatusNACK, result, ref.start)
	ref.m.Nack()
	return true
}

func (ref *pubsubWrapper) ConcurrencyKey() string {
	return ref.m.Attributes["concurrency_key"]
}
