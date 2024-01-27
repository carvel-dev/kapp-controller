// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

// Package metrics to define all prometheus metric methods
package metrics

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// ReconcileTimeMetrics holds reconcile time metrics
type ReconcileTimeMetrics struct {
	reconcileTimeSeconds         *prometheus.GaugeVec
	reconcileDeployTimeSeconds   *prometheus.GaugeVec
	reconcileFetchTimeSeconds    *prometheus.GaugeVec
	reconcileTemplateTimeSeconds *prometheus.GaugeVec
}

var (
	timeMetricsOnce sync.Once
)

// NewReconcileTimeMetrics creates ReconcileTimeMetrics object
func NewReconcileTimeMetrics() *ReconcileTimeMetrics {
	const (
		metricNamespace     = "kappctrl"
		resourceTypeLabel   = "controller"
		resourceNameLabel   = "name"
		firstReconcileLabel = "firstReconcile"
		namespaceLabel      = "namespace"
	)
	return &ReconcileTimeMetrics{
		reconcileTimeSeconds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricNamespace,
				Name:      "reconcile_time_seconds",
				Help:      "Overall time taken to reconcile a CR",
			},
			[]string{resourceTypeLabel, resourceNameLabel, namespaceLabel, firstReconcileLabel},
		),
		reconcileFetchTimeSeconds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricNamespace,
				Name:      "reconcile_fetch_time_seconds",
				Help:      "Time taken to perform a fetch for a CR",
			},
			[]string{resourceTypeLabel, resourceNameLabel, namespaceLabel, firstReconcileLabel},
		),
		reconcileTemplateTimeSeconds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricNamespace,
				Name:      "reconcile_template_time_seconds",
				Help:      "Time taken to perform a templating for a CR",
			},
			[]string{resourceTypeLabel, resourceNameLabel, namespaceLabel, firstReconcileLabel},
		),
		reconcileDeployTimeSeconds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricNamespace,
				Name:      "reconcile_deploy_time_seconds",
				Help:      "Time taken to perform a deploy for a CR",
			},
			[]string{resourceTypeLabel, resourceNameLabel, namespaceLabel, firstReconcileLabel},
		),
	}
}

// RegisterAllMetrics registers reconcile time prometheus metrics.
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

// RegisterOverallTime sets overall time
func (tm *ReconcileTimeMetrics) RegisterOverallTime(resourceType, name, namespace string, firstReconcile bool, time time.Duration) {
	tm.reconcileTimeSeconds.WithLabelValues(resourceType, name, namespace, strconv.FormatBool(firstReconcile)).Set(time.Seconds())
}

// RegisterFetchTime sets fetch time
func (tm *ReconcileTimeMetrics) RegisterFetchTime(resourceType, name, namespace string, firstReconcile bool, time time.Duration) {
	tm.reconcileFetchTimeSeconds.WithLabelValues(resourceType, name, namespace, strconv.FormatBool(firstReconcile)).Set(time.Seconds())
}

// RegisterTemplateTime sets template time
func (tm *ReconcileTimeMetrics) RegisterTemplateTime(resourceType, name, namespace string, firstReconcile bool, time time.Duration) {
	tm.reconcileTemplateTimeSeconds.WithLabelValues(resourceType, name, namespace, strconv.FormatBool(firstReconcile)).Set(time.Seconds())
}

// RegisterDeployTime sets deploy time
func (tm *ReconcileTimeMetrics) RegisterDeployTime(resourceType, name, namespace string, firstReconcile bool, time time.Duration) {
	tm.reconcileDeployTimeSeconds.WithLabelValues(resourceType, name, namespace, strconv.FormatBool(firstReconcile)).Set(time.Seconds())
}
