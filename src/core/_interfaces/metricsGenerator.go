package _interfaces

import "time"

type MetricsGenerator interface {
	EventReceived(publishedTime time.Time)
	EventProcessed(handshake string, result string, requestReceived time.Time)
}
