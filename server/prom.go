package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

var proxyLatency = prometheus.NewHistogramVec( //nolint:gochecknoglobals prom
	prometheus.HistogramOpts{
		Name: "circa_target_duration_seconds",
		Help: "Histogram for requests to target",
	},
	[]string{"target", "method", "route", "status"},
)
var requestsLatency = prometheus.NewHistogramVec( //nolint:gochecknoglobals prom
	prometheus.HistogramOpts{
		Name: "circa_request_duration_seconds",
		Help: "Histogram income requests ",
	},
	[]string{"method", "route", "status"},
)

func RegisterMetrics() {
	prometheus.MustRegister(proxyLatency)
	prometheus.MustRegister(requestsLatency)
}
