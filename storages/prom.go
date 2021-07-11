package storages

import (
	"github.com/prometheus/client_golang/prometheus"
)

var operationHistogram = prometheus.NewHistogramVec( //nolint:gochecknoglobals prom
	prometheus.HistogramOpts{
		Name: "circa_storage_operation_duration_seconds",
		Help: "Histogram for operations in storages",
	},
	[]string{"storage", "operation"},
)

func RegisterMetrics() {
	prometheus.MustRegister(operationHistogram)
}
