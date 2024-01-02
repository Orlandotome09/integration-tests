package metrics

import (
	_ "github.com/prometheus/client_golang/prometheus/push"
	"time"
)

//HTTP application
type HTTP struct {
	Handler    string
	Method     string
	StatusCode string
	Partner    string
	RequestId  string
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   float64
}

//NewHTTP create a new HTTP app
func NewHTTP(handler, method, requestId string) *HTTP {
	return &HTTP{
		Handler:   handler,
		Method:    method,
		RequestId: requestId,
	}
}

//Started start monitoring the app
func (h *HTTP) Started() {
	h.StartedAt = time.Now()
}

// Finished app finished
func (h *HTTP) Finished() {
	h.FinishedAt = time.Now()
	h.Duration = time.Since(h.StartedAt).Seconds()
}
