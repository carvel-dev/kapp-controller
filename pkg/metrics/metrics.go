// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// ServerMetrics holds server metrics
type ServerMetrics struct {
	reconcileAttemptTotal       *prometheus.CounterVec
	reconcileSuccessTotal       *prometheus.CounterVec
	reconcileFailureTotal       *prometheus.CounterVec
	reconcileDeleteAttemptTotal *prometheus.CounterVec
	reconcileDeleteFailedTotal  *prometheus.CounterVec
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
	kappNamespaceLabel = "namespace"
)

// NewServerMetrics creates ServerMetrics object
func NewServerMetrics() *ServerMetrics {
	return &ServerMetrics{
		reconcileAttemptTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      reconcileAttemptTotal,
				Help:      "Total number of attempted reconciles",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		reconcileSuccessTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      reconcileSuccessTotal,
				Help:      "Total number of success reconciles",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		reconcileFailureTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      reconcileFailureTotal,
				Help:      "Total number of failed reconciles",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		reconcileDeleteAttemptTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      reconcileDeleteAttemptTotal,
				Help:      "Total number of attempted reconcile deletion",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
		reconcileDeleteFailedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricNamespace,
				Name:      reconcileDeleteFailedTotal,
				Help:      "Total number of failed reconcile deletion",
			},
			[]string{kappNameLabel, kappNamespaceLabel},
		),
	}
}

// RegisterAllMetrics registers all prometheus metrics.
func (sm *ServerMetrics) RegisterAllMetrics() {
	metrics.Registry.MustRegister(
		sm.reconcileAttemptTotal,
		sm.reconcileSuccessTotal,
		sm.reconcileFailureTotal,
		sm.reconcileDeleteAttemptTotal,
		sm.reconcileDeleteFailedTotal,
	)
}

// InitMetrics initializes metrics
func (sm *ServerMetrics) InitMetrics(appName string, namespace string) {
	// Initializes counter metrics
	sm.reconcileAttemptTotal.WithLabelValues(appName, namespace).Add(0)
	sm.reconcileSuccessTotal.WithLabelValues(appName, namespace).Add(0)
	sm.reconcileFailureTotal.WithLabelValues(appName, namespace).Add(0)
	sm.reconcileDeleteAttemptTotal.WithLabelValues(appName, namespace).Add(0)
	sm.reconcileDeleteFailedTotal.WithLabelValues(appName, namespace).Add(0)
}

// DeleteMetrics deletes metrics
func (sm *ServerMetrics) DeleteMetrics(appName string, namespace string) {
	// Delete counter metrics
	sm.reconcileAttemptTotal.DeleteLabelValues(appName, namespace)
	sm.reconcileSuccessTotal.DeleteLabelValues(appName, namespace)
	sm.reconcileFailureTotal.DeleteLabelValues(appName, namespace)
	sm.reconcileDeleteAttemptTotal.DeleteLabelValues(appName, namespace)
	sm.reconcileDeleteFailedTotal.DeleteLabelValues(appName, namespace)
}

// RegisterReconcileAttempt increments reconcileAttemptTotal
func (sm *ServerMetrics) RegisterReconcileAttempt(appName string, namespace string) {
	sm.reconcileAttemptTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileSuccess increments reconcileSuccessTotal
func (sm *ServerMetrics) RegisterReconcileSuccess(appName string, namespace string) {
	sm.reconcileSuccessTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileFailure increments reconcileFailureTotal
func (sm *ServerMetrics) RegisterReconcileFailure(appName string, namespace string) {
	sm.reconcileFailureTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileDeleteAttempt increments reconcileDeleteAttemptTotal
func (sm *ServerMetrics) RegisterReconcileDeleteAttempt(appName string, namespace string) {
	sm.reconcileDeleteAttemptTotal.WithLabelValues(appName, namespace).Inc()
}

// RegisterReconcileDeleteFailed increments reconcileDeleteFailedTotal
func (sm *ServerMetrics) RegisterReconcileDeleteFailed(appName string, namespace string) {
	sm.reconcileDeleteFailedTotal.WithLabelValues(appName, namespace).Inc()
}
