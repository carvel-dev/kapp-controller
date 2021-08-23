package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type ServerMetrics struct {
	allMetrics map[string]prometheus.Collector
}

const (
	metricNamespace = "kapp"

	// kpp metrics
	reconcileAttemptTotal       = "reconcile_attempt_total"
	reconcileSuccessTotal       = "reconcile_success_total"
	reconcileFailureTotal       = "reconcile_failure_total"
	reconcileDeleteAttemptTotal = "reconcile_delete_attempt_total"
	reconcileDeleteFailedTotal  = "reconcile_delete_failed_total"

	// Labels
	kappNameLabel      = "appName"
	KappNamespaceLabel = "namespace"
)

var (
	allMetricsMap = ServerMetrics{
		allMetrics: map[string]prometheus.Collector{
			reconcileAttemptTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: metricNamespace,
					Name:      reconcileAttemptTotal,
					Help:      "Total number of attempted reconciles",
				},
				[]string{kappNameLabel, KappNamespaceLabel},
			),
			reconcileSuccessTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: metricNamespace,
					Name:      reconcileSuccessTotal,
					Help:      "Total number of success reconciles",
				},
				[]string{kappNameLabel, KappNamespaceLabel},
			),
			reconcileFailureTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: metricNamespace,
					Name:      reconcileFailureTotal,
					Help:      "Total number of failed reconciles",
				},
				[]string{kappNameLabel, KappNamespaceLabel},
			),
			reconcileDeleteAttemptTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: metricNamespace,
					Name:      reconcileDeleteAttemptTotal,
					Help:      "Total number of attempted reconcile deletion",
				},
				[]string{kappNameLabel, KappNamespaceLabel},
			),
			reconcileDeleteFailedTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: metricNamespace,
					Name:      reconcileDeleteFailedTotal,
					Help:      "Total number of failed reconcile deletion",
				},
				[]string{kappNameLabel, KappNamespaceLabel},
			),
		},
	}
)

// RegisterAllMetrics registers all prometheus metrics.
func init() {
	for _, pm := range allMetricsMap.allMetrics {
		metrics.Registry.MustRegister(pm)
	}
}

// InitMetrics initializes counter metrics .
func InitMetrics(appName string, namespace string) {
	for key, _ := range allMetricsMap.allMetrics {
		if c, ok := allMetricsMap.allMetrics[key].(*prometheus.CounterVec); ok {
			c.WithLabelValues(appName, namespace).Add(0)
		}
	}
}

// DeleteMetrics initializes counter metrics .
func DeleteMetrics(appName string, namespace string) {
	for key, _ := range allMetricsMap.allMetrics {
		if c, ok := allMetricsMap.allMetrics[key].(*prometheus.CounterVec); ok {
			c.DeleteLabelValues(appName, namespace)
		}
	}
}

func RegisterReconcileAttempt(appName string, namespace string) {
	if c, ok := allMetricsMap.allMetrics[reconcileAttemptTotal].(*prometheus.CounterVec); ok {
		c.WithLabelValues(appName, namespace).Inc()
	}
}

func RegisterReconcileSuccess(appName string, namespace string) {
	if c, ok := allMetricsMap.allMetrics[reconcileSuccessTotal].(*prometheus.CounterVec); ok {
		c.WithLabelValues(appName, namespace).Inc()
	}
}

func RegisterReconcileFailure(appName string, namespace string) {
	if c, ok := allMetricsMap.allMetrics[reconcileFailureTotal].(*prometheus.CounterVec); ok {
		c.WithLabelValues(appName, namespace).Inc()
	}
}

func RegisterReconcileDeleteAttempt(appName string, namespace string) {
	if c, ok := allMetricsMap.allMetrics[reconcileDeleteAttemptTotal].(*prometheus.CounterVec); ok {
		c.WithLabelValues(appName, namespace).Inc()
	}
}

func RegisterReconcileDeleteFailed(appName string, namespace string) {
	if c, ok := allMetricsMap.allMetrics[reconcileDeleteFailedTotal].(*prometheus.CounterVec); ok {
		c.WithLabelValues(appName, namespace).Inc()
	}
}
