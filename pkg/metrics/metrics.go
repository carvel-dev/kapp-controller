// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package metrics

// Metrics holds all metrics
type Metrics struct {
	*ReconcileCountMetrics
	*ReconcileTimeMetrics
	IsFirstReconcile bool
}

// NewMetrics is a factory function that returns a new instance of Metrics.
func NewMetrics() *Metrics {
	return &Metrics{
		ReconcileCountMetrics: NewCountMetrics(),
		ReconcileTimeMetrics:  NewReconcileTimeMetrics(),
		IsFirstReconcile:      false,
	}
}
