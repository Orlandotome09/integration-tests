package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func ExposeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":7777", nil)
}
