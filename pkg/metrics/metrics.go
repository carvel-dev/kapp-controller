package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

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
	reconcileAttemptTotalKey = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Name:      reconcileAttemptTotal,
		Help:      "Total number of attempted reconciles",
	}, []string{kappNameLabel, KappNamespaceLabel},
	)
	reconcileSuccessTotalKey = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Name:      reconcileSuccessTotal,
		Help:      "Total number of success reconciles",
	}, []string{kappNameLabel, KappNamespaceLabel},
	)
	reconcileFailureTotalKey = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Name:      reconcileFailureTotal,
		Help:      "Total number of failed reconciles",
	}, []string{kappNameLabel, KappNamespaceLabel},
	)
	reconcileDeleteAttemptTotalKey = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Name:      reconcileDeleteAttemptTotal,
		Help:      "Total number of attempted reconcile deletion",
	}, []string{kappNameLabel, KappNamespaceLabel},
	)
	reconcileDeleteFailedTotalKey = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Name:      reconcileDeleteFailedTotal,
		Help:      "Total number of failed reconcile deletion",
	}, []string{kappNameLabel, KappNamespaceLabel},
	)
)

// RegisterAllMetrics registers all prometheus metrics.
func init() {
	metrics.Registry.MustRegister(reconcileAttemptTotalKey)
	metrics.Registry.MustRegister(reconcileSuccessTotalKey)
	metrics.Registry.MustRegister(reconcileFailureTotalKey)
	metrics.Registry.MustRegister(reconcileDeleteAttemptTotalKey)
	metrics.Registry.MustRegister(reconcileDeleteFailedTotalKey)
}

// InitMetrics initializes counter metrics .
func InitMetrics(appName string, namespace string) {
	reconcileAttemptTotalKey.WithLabelValues(appName, namespace).Add(0)
	reconcileSuccessTotalKey.WithLabelValues(appName, namespace).Add(0)
	reconcileFailureTotalKey.WithLabelValues(appName, namespace).Add(0)
	reconcileDeleteAttemptTotalKey.WithLabelValues(appName, namespace).Add(0)
	reconcileDeleteFailedTotalKey.WithLabelValues(appName, namespace).Add(0)
}

// DeleteMetrics initializes counter metrics .
func DeleteMetrics(appName string, namespace string) {
	reconcileAttemptTotalKey.DeleteLabelValues(appName, namespace)
	reconcileSuccessTotalKey.DeleteLabelValues(appName, namespace)
	reconcileFailureTotalKey.DeleteLabelValues(appName, namespace)
	reconcileDeleteAttemptTotalKey.DeleteLabelValues(appName, namespace)
	reconcileDeleteFailedTotalKey.DeleteLabelValues(appName, namespace)
}

func RegisterReconcileAttempt(appName string, namespace string) {
	reconcileAttemptTotalKey.WithLabelValues(appName, namespace).Inc()
}

func RegisterReconcileSuccess(appName string, namespace string) {
	reconcileSuccessTotalKey.WithLabelValues(appName, namespace).Inc()
}

func RegisterReconcileFailure(appName string, namespace string) {
	reconcileFailureTotalKey.WithLabelValues(appName, namespace).Inc()
}

func RegisterReconcileDeleteAttempt(appName string, namespace string) {
	reconcileDeleteAttemptTotalKey.WithLabelValues(appName, namespace).Inc()
}

func RegisterReconcileDeleteFailed(appName string, namespace string) {
	reconcileDeleteFailedTotalKey.WithLabelValues(appName, namespace).Inc()
}
