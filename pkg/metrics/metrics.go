// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package metrics

// Metrics holds all metrics
type Metrics struct {
	*ReconcileCountMetrics
	*ReconcileTimeMetrics
	IsFirstReconcile bool
}
