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
	reconcileAttemptTotal       *prometheus.CounterVec
	reconcileSuccessTotal       *prometheus.CounterVec
	reconcileFailureTotal       *prometheus.CounterVec
	reconcileDeleteAttemptTotal *prometheus.CounterVec
	reconcileDeleteFailedTotal  *prometheus.CounterVec
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
		reconcileAttemptTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_attempt_total",
				Help:      "Total number of attempted reconciles",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		reconcileSuccessTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_success_total",
				Help:      "Total number of succeeded reconciles",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		reconcileFailureTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_failure_total",
				Help:      "Total number of failed reconciles",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		reconcileDeleteAttemptTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      "app_reconcile_delete_attempt_total",
				Help:      "Total number of attempted reconcile deletion",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		reconcileDeleteFailedTotal: prometheus.NewCounterVec(
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
			am.reconcileAttemptTotal,
			am.reconcileSuccessTotal,
			am.reconcileFailureTotal,
			am.reconcileDeleteAttemptTotal,
			am.reconcileDeleteFailedTotal,
		)
	})
}

// InitMetrics initializes metrics
func (am *AppMetrics) InitMetrics(appName string, namespace string) {
	// Initializes counter metrics
	am.reconcileAttemptTotal.WithLabelValues(appName, namespace).Add(0)
	am.reconcileSuccessTotal.WithLabelValues(appName, namespace).Add(0)
	am.reconcileFailureTotal.WithLabelValues(appName, namespace).Add(0)
	am.reconcileDeleteAttemptTotal.WithLabelValues(appName, namespace).Add(0)
	am.reconcileDeleteFailedTotal.WithLabelValues(appName, namespace).Add(0)
}

// DeleteMetrics deletes metrics
func (am *AppMetrics) DeleteMetrics(appName string, namespace string) {
	// Delete counter metrics
	am.reconcileAttemptTotal.DeleteLabelValues(appName, namespace)
	am.reconcileSuccessTotal.DeleteLabelValues(appName, namespace)
	am.reconcileFailureTotal.DeleteLabelValues(appName, namespace)
	am.reconcileDeleteAttemptTotal.DeleteLabelValues(appName, namespace)
	am.reconcileDeleteFailedTotal.DeleteLabelValues(appName, namespace)
}

// RegisterReconcileAttempt increments reconcileAttemptTotal
func (am *AppMetrics) RegisterReconcileAttempt(appName string, namespace string) {
	am.reconcileAttemptTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileSuccess increments reconcileSuccessTotal
func (am *AppMetrics) RegisterReconcileSuccess(appName string, namespace string) {
	am.reconcileSuccessTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileFailure increments reconcileFailureTotal
func (am *AppMetrics) RegisterReconcileFailure(appName string, namespace string) {
	am.reconcileFailureTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileDeleteAttempt increments reconcileDeleteAttemptTotal
func (am *AppMetrics) RegisterReconcileDeleteAttempt(appName string, namespace string) {
	am.reconcileDeleteAttemptTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileDeleteFailed increments reconcileDeleteFailedTotal
func (am *AppMetrics) RegisterReconcileDeleteFailed(appName string, namespace string) {
	am.reconcileDeleteFailedTotal.WithLabelValues(appName, namespace).Inc()
}
