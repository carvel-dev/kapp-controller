// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// AppMetrics holds server metrics
type AppMetrics struct {
	ReconcileAttemptTotal       *prometheus.CounterVec
	ReconcileSuccessTotal       *prometheus.CounterVec
	ReconcileFailureTotal       *prometheus.CounterVec
	ReconcileDeleteAttemptTotal *prometheus.CounterVec
	ReconcileDeleteFailedTotal  *prometheus.CounterVec
}

var (
	once sync.Once
)

// NewAppMetrics creates AppMetrics object
func NewAppMetrics() *AppMetrics {
	const (
		metricNamespace    = "kappctrl"
		kappNameLabel      = "app_name"
		kappNamespaceLabel = "namespace"
	)
	return &AppMetrics{
		ReconcileAttemptTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_attempt_total",
				Help:      "Total number of attempted reconciles",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		ReconcileSuccessTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_success_total",
				Help:      "Total number of succeeded reconciles",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		ReconcileFailureTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_failure_total",
				Help:      "Total number of failed reconciles",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		ReconcileDeleteAttemptTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_delete_attempt_total",
				Help:      "Total number of attempted reconcile deletion",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		ReconcileDeleteFailedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_delete_failed_total",
				Help:      "Total number of failed reconcile deletion",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
	}
}

// RegisterAllMetrics registers all prometheus metrics.
func (am *AppMetrics) RegisterAllMetrics() {
	once.Do(func() {
		metrics.Registry.MustRegister(
			am.ReconcileAttemptTotal,
			am.ReconcileSuccessTotal,
			am.ReconcileFailureTotal,
			am.ReconcileDeleteAttemptTotal,
			am.ReconcileDeleteFailedTotal,
		)
	})
}

// InitMetrics initializes metrics
func (am *AppMetrics) InitMetrics(appName string, namespace string) {
	// Initializes counter metrics
	am.ReconcileAttemptTotal.WithLabelValues(appName, namespace).Add(0)
	am.ReconcileSuccessTotal.WithLabelValues(appName, namespace).Add(0)
	am.ReconcileFailureTotal.WithLabelValues(appName, namespace).Add(0)
	am.ReconcileDeleteAttemptTotal.WithLabelValues(appName, namespace).Add(0)
	am.ReconcileDeleteFailedTotal.WithLabelValues(appName, namespace).Add(0)
}

// DeleteMetrics deletes metrics
func (am *AppMetrics) DeleteMetrics(appName string, namespace string) {
	// Delete counter metrics
	am.ReconcileAttemptTotal.DeleteLabelValues(appName, namespace)
	am.ReconcileSuccessTotal.DeleteLabelValues(appName, namespace)
	am.ReconcileFailureTotal.DeleteLabelValues(appName, namespace)
	am.ReconcileDeleteAttemptTotal.DeleteLabelValues(appName, namespace)
	am.ReconcileDeleteFailedTotal.DeleteLabelValues(appName, namespace)
}

// RegisterReconcileAttempt increments reconcileAttemptTotal
func (am *AppMetrics) RegisterReconcileAttempt(appName string, namespace string) {
	am.ReconcileAttemptTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileSuccess increments reconcileSuccessTotal
func (am *AppMetrics) RegisterReconcileSuccess(appName string, namespace string) {
	am.ReconcileSuccessTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileFailure increments reconcileFailureTotal
func (am *AppMetrics) RegisterReconcileFailure(appName string, namespace string) {
	am.ReconcileFailureTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileDeleteAttempt increments reconcileDeleteAttemptTotal
func (am *AppMetrics) RegisterReconcileDeleteAttempt(appName string, namespace string) {
	am.ReconcileDeleteAttemptTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileDeleteFailed increments reconcileDeleteFailedTotal
func (am *AppMetrics) RegisterReconcileDeleteFailed(appName string, namespace string) {
	am.ReconcileDeleteFailedTotal.WithLabelValues(appName, namespace).Inc()
}
