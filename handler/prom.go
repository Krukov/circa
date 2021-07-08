package handler

import (
	"github.com/prometheus/client_golang/prometheus"
)

var routeHandlerCount = prometheus.NewCounterVec( //nolint:gochecknoglobals prom
	prometheus.CounterOpts{
		Name:        "circa_handle_count",
		Help:        "Handle counter",
	},
	[]string{"rule", "route", "key", "status"},
)

var handlersGauge = prometheus.NewGaugeVec( //nolint:gochecknoglobals prom
	prometheus.GaugeOpts{
		Name:        "circa_handlers_gauge",
		Help:        "Handlers config by routes",
	},
	[]string{"rule", "route"},
)


func RegisterMetrics() {
	prometheus.MustRegister(routeHandlerCount)
	prometheus.MustRegister(handlersGauge)
}
