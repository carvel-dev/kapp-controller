// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type serverMetrics struct {
	metrics map[string]prometheus.Collector
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

var (
	sm = serverMetrics{
		metrics: map[string]prometheus.Collector{
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
		},
	}
)

// RegisterAllMetrics registers all prometheus metrics.
func init() {
	for _, pm := range sm.metrics {
		metrics.Registry.MustRegister(pm)
	}
}

// InitMetrics initializes counter metrics
func InitMetrics(appName string, namespace string) {
	for key := range sm.metrics {
		if c, ok := sm.metrics[key].(*prometheus.CounterVec); ok {
			c.WithLabelValues(appName, namespace).Add(0)
		}
	}
}

// DeleteMetrics deletes counter metrics
func DeleteMetrics(appName string, namespace string) {
	for key := range sm.metrics {
		if c, ok := sm.metrics[key].(*prometheus.CounterVec); ok {
			c.DeleteLabelValues(appName, namespace)
		}
	}
}

// RegisterReconcileAttempt increments reconcileAttemptTotal
func RegisterReconcileAttempt(appName string, namespace string) {
	if c, ok := sm.metrics[reconcileAttemptTotal].(*prometheus.CounterVec); ok {
		c.WithLabelValues(appName, namespace).Inc()
	}
}

// RegisterReconcileSuccess increments reconcileSuccessTotal
func RegisterReconcileSuccess(appName string, namespace string) {
	if c, ok := sm.metrics[reconcileSuccessTotal].(*prometheus.CounterVec); ok {
		c.WithLabelValues(appName, namespace).Inc()
	}
}

// RegisterReconcileFailure increments reconcileFailureTotal
func RegisterReconcileFailure(appName string, namespace string) {
	if c, ok := sm.metrics[reconcileFailureTotal].(*prometheus.CounterVec); ok {
		c.WithLabelValues(appName, namespace).Inc()
	}
}

// RegisterReconcileDeleteAttempt increments reconcileDeleteAttemptTotal
func RegisterReconcileDeleteAttempt(appName string, namespace string) {
	if c, ok := sm.metrics[reconcileDeleteAttemptTotal].(*prometheus.CounterVec); ok {
		c.WithLabelValues(appName, namespace).Inc()
	}
}

// RegisterReconcileDeleteFailed increments reconcileDeleteFailedTotal
func RegisterReconcileDeleteFailed(appName string, namespace string) {
	if c, ok := sm.metrics[reconcileDeleteFailedTotal].(*prometheus.CounterVec); ok {
		c.WithLabelValues(appName, namespace).Inc()
	}
}
