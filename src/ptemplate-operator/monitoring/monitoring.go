package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	PtemplateCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "Ptemplate_reconcile_count",
			Help: "Total number of times the operartor has run reconcile().",
		},
	)

	PtemplateGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "Ptemplate_consumers",
			Help: "Number of existing Ptemplate consumers",
		},
	)
)

// RegisterMetrics will register metrics with the global prometheus registry
func RegisterMetrics() {
	metrics.Registry.MustRegister(PtemplateCounter)
	metrics.Registry.MustRegister(PtemplateGauge)
}
