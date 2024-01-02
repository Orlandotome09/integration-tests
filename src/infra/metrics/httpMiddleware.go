package metrics

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpConcurrentProcess.Inc()
		appMetric := NewHTTP(c.FullPath(), c.Request.Method, "")
		appMetric.Started()
		c.Next()
		statusCode := c.Writer.Status()
		appMetric.Finished()
		appMetric.StatusCode = strconv.Itoa(statusCode)
		sendMetrics(appMetric)
		httpConcurrentProcess.Dec()
	}
}

// Send metrics to prometheus server
func sendMetrics(h *HTTP) {
	httpDuration.WithLabelValues(h.Handler, h.Method, h.StatusCode).Observe(h.Duration)
	httpRequests.WithLabelValues(h.Handler, h.Method, h.StatusCode).Inc()
	httpResponseStatus.WithLabelValues(h.Handler, h.Method, h.StatusCode).Inc()

}
