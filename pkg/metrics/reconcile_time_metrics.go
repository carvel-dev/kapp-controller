package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"sync"
	"time"
)

type ReconcileTimeMetrics struct {
	reconcileTimeSeconds         *prometheus.GaugeVec
	reconcileDeployTimeSeconds   *prometheus.GaugeVec
	reconcileFetchTimeSeconds    *prometheus.GaugeVec
	reconcileTemplateTimeSeconds *prometheus.GaugeVec
}

var (
	timeMetricsOnce sync.Once
)

func NewReconcileTimeMetrics() *ReconcileTimeMetrics {
	const (
		metricNamespace     = "kappctrl_reconcile_time_seconds"
		resourceTypeLabel   = "controller"
		resourceNameLabel   = "name"
		firstReconcileLabel = "firstReconcile"
		namespace           = "namespace"
	)
	return &ReconcileTimeMetrics{
		reconcileTimeSeconds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricNamespace,
				Name:      "reconcile_time_seconds",
				Help:      "Overall time taken to reconcile a CR",
			},
			[]string{resourceTypeLabel, resourceNameLabel, namespace, firstReconcileLabel},
		),
		reconcileFetchTimeSeconds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricNamespace,
				Name:      "reconcile_fetch_time_seconds",
				Help:      "Time taken to perform a fetch for a CR",
			},
			[]string{resourceTypeLabel, resourceNameLabel, namespace, firstReconcileLabel},
		),
		reconcileTemplateTimeSeconds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricNamespace,
				Name:      "reconcile_template_time_seconds",
				Help:      "Time taken to perform a templating for a CR",
			},
			[]string{resourceTypeLabel, resourceNameLabel, namespace, firstReconcileLabel},
		),
		reconcileDeployTimeSeconds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricNamespace,
				Name:      "reconcile_deploy_time_seconds",
				Help:      "Time taken to perform a deploy for a CR",
			},
			[]string{resourceTypeLabel, resourceNameLabel, namespace, firstReconcileLabel},
		),
	}
}

func (tm *ReconcileTimeMetrics) RegisterAllMetrics() {
	timeMetricsOnce.Do(func() {
		metrics.Registry.MustRegister(
			tm.reconcileTimeSeconds,
			tm.reconcileFetchTimeSeconds,
			tm.reconcileTemplateTimeSeconds,
			tm.reconcileDeployTimeSeconds,
		)
	})
}

func (tm *ReconcileTimeMetrics) RegisterOverallTime(resourceType, name, namespace, firstReconcile string, time time.Duration) {
	tm.reconcileTimeSeconds.WithLabelValues(resourceType, name, namespace, firstReconcile).Set(time.Seconds())
}

func (tm *ReconcileTimeMetrics) RegisterFetchTime(resourceType, name, namespace, firstReconcile string, time time.Duration) {
	tm.reconcileFetchTimeSeconds.WithLabelValues(resourceType, name, namespace, firstReconcile).Set(time.Seconds())
}

func (tm *ReconcileTimeMetrics) RegisterTemplateTime(resourceType, name, namespace, firstReconcile string, time time.Duration) {
	tm.reconcileTemplateTimeSeconds.WithLabelValues(resourceType, name, namespace, firstReconcile).Set(time.Seconds())
}

func (tm *ReconcileTimeMetrics) RegisterDeployTime(resourceType, name, namespace, firstReconcile string, time time.Duration) {
	tm.reconcileDeployTimeSeconds.WithLabelValues(resourceType, name, namespace, firstReconcile).Set(time.Seconds())
}
