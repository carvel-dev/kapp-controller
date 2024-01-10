// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// ReconcileCountMetrics holds server metrics
type ReconcileCountMetrics struct {
	reconcileAttemptTotal       *prometheus.CounterVec
	reconcileSuccessTotal       *prometheus.CounterVec
	reconcileFailureTotal       *prometheus.CounterVec
	reconcileDeleteAttemptTotal *prometheus.CounterVec
	reconcileDeleteFailedTotal  *prometheus.CounterVec
}

var (
	once sync.Once
)

// NewCountMetrics creates ReconcileCountMetrics object
func NewCountMetrics() *ReconcileCountMetrics {
	const (
		metricNamespace    = "kappctrl"
		kappNameLabel      = "name"
		kappNamespaceLabel = "namespace"
		resourceTypeLabel  = "controller"
	)
	return &ReconcileCountMetrics{
		reconcileAttemptTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_attempt_total",
				Help:      "Total number of attempted reconciles",
			},
			[]string{resourceTypeLabel, kappNameLabel, kappNamespaceLabel},
		),
		reconcileSuccessTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_success_total",
				Help:      "Total number of succeeded reconciles",
			},
			[]string{resourceTypeLabel, kappNameLabel, kappNamespaceLabel},
		),
		reconcileFailureTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_failure_total",
				Help:      "Total number of failed reconciles",
			},
			[]string{resourceTypeLabel, kappNameLabel, kappNamespaceLabel},
		),
		reconcileDeleteAttemptTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_delete_attempt_total",
				Help:      "Total number of attempted reconcile deletions",
			},
			[]string{resourceTypeLabel, kappNameLabel, kappNamespaceLabel},
		),
		reconcileDeleteFailedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_delete_failed_total",
				Help:      "Total number of failed reconcile deletions",
			},
			[]string{resourceTypeLabel, kappNameLabel, kappNamespaceLabel},
		),
	}
}

// RegisterAllMetrics registers all prometheus metrics.
func (am *ReconcileCountMetrics) RegisterAllMetrics() {
	once.Do(func() {
		metrics.Registry.MustRegister(
			am.reconcileAttemptTotal,
			am.reconcileSuccessTotal,
			am.reconcileFailureTotal,
			am.reconcileDeleteAttemptTotal,
			am.reconcileDeleteFailedTotal,
		)
	})
}

// InitMetrics initializes metrics
func (am *ReconcileCountMetrics) InitMetrics(resourceType, name, namespace string) {
	// Initializes counter metrics
	am.reconcileAttemptTotal.WithLabelValues(resourceType, name, namespace).Add(0)
	am.reconcileSuccessTotal.WithLabelValues(resourceType, name, namespace).Add(0)
	am.reconcileFailureTotal.WithLabelValues(resourceType, name, namespace).Add(0)
	am.reconcileDeleteAttemptTotal.WithLabelValues(resourceType, name, namespace).Add(0)
	am.reconcileDeleteFailedTotal.WithLabelValues(resourceType, name, namespace).Add(0)
}

// DeleteMetrics deletes metrics
func (am *ReconcileCountMetrics) DeleteMetrics(resourceType, name, namespace string) {
	// Delete counter metrics
	am.reconcileAttemptTotal.DeleteLabelValues(resourceType, name, namespace)
	am.reconcileSuccessTotal.DeleteLabelValues(resourceType, name, namespace)
	am.reconcileFailureTotal.DeleteLabelValues(resourceType, name, namespace)
	am.reconcileDeleteAttemptTotal.DeleteLabelValues(resourceType, name, namespace)
	am.reconcileDeleteFailedTotal.DeleteLabelValues(resourceType, name, namespace)
}

// RegisterReconcileAttempt increments reconcileAttemptTotal
func (am *ReconcileCountMetrics) RegisterReconcileAttempt(resourceType, appName, namespace string) {
	am.reconcileAttemptTotal.WithLabelValues(resourceType, appName, namespace).Inc()
}

// RegisterReconcileSuccess increments reconcileSuccessTotal
func (am *ReconcileCountMetrics) RegisterReconcileSuccess(resourceType, appName, namespace string) {
	am.reconcileSuccessTotal.WithLabelValues(resourceType, appName, namespace).Inc()
}

// RegisterReconcileFailure increments reconcileFailureTotal
func (am *ReconcileCountMetrics) RegisterReconcileFailure(resourceType, appName, namespace string) {
	am.reconcileFailureTotal.WithLabelValues(resourceType, appName, namespace).Inc()
}

// RegisterReconcileDeleteAttempt increments reconcileDeleteAttemptTotal
func (am *ReconcileCountMetrics) RegisterReconcileDeleteAttempt(resourceType, appName, namespace string) {
	am.reconcileDeleteAttemptTotal.WithLabelValues(resourceType, appName, namespace).Inc()
}

// RegisterReconcileDeleteFailed increments reconcileDeleteFailedTotal
func (am *ReconcileCountMetrics) RegisterReconcileDeleteFailed(resourceType, appName, namespace string) {
	am.reconcileDeleteFailedTotal.WithLabelValues(resourceType, appName, namespace).Inc()
}

// GetReconcileAttemptCounterValue return reconcile count
func (am *ReconcileCountMetrics) GetReconcileAttemptCounterValue(resourceType, appName, namespace string) int64 {
	var m = &dto.Metric{}
	if err := am.reconcileAttemptTotal.WithLabelValues(resourceType, appName, namespace).Write(m); err != nil {
		return 0
	}
	return int64(m.Counter.GetValue())
}
