// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"math"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type ReconcileTimer struct {
	app v1alpha1.App
}

func NewReconcileTimer(app v1alpha1.App) ReconcileTimer {
	return ReconcileTimer{*app.DeepCopy()}
}

func (rt ReconcileTimer) DurationUntilReady(err error) time.Duration {
	if err != nil || rt.hasReconcileStatus(v1alpha1.ReconcileFailed) {
		return rt.failureSyncPeriod()
	}

	return rt.applyJitter(rt.syncPeriod())
}

func (rt ReconcileTimer) IsReadyAt(timeAt time.Time) bool {
	// Did resource spec change?
	if rt.app.Status.ObservedGeneration != rt.app.Generation {
		return true
	}

	// TODO: Is this needed due to first case statement?
	// If paused, then no reconcilation until unpaused
	if rt.app.Spec.Paused {
		return false
	}

	var lastReconcileTime time.Time

	// Use latest deploy time if available, otherwise fallback to fetch
	// If no timestamp is available, enqueue immediately
	lastDeploy := rt.app.Status.Deploy
	lastFetch := rt.app.Status.Fetch
	if lastDeploy != nil && !lastDeploy.UpdatedAt.Time.IsZero() {
		lastReconcileTime = rt.app.Status.Deploy.UpdatedAt.Time
	} else {
		if lastFetch == nil {
			return true
		}
		lastReconcileTime = lastFetch.UpdatedAt.Time
	}

	if rt.hasReconcileStatus(v1alpha1.ReconcileFailed) {
		if timeAt.UTC().Sub(lastReconcileTime) >= rt.failureSyncPeriod() {
			return true
		}
	}

	// Did we deploy too long ago?
	if timeAt.UTC().Sub(lastReconcileTime) >= rt.syncPeriod() {
		return true
	}

	return false
}

func (rt ReconcileTimer) syncPeriod() time.Duration {
	const defaultSyncPeriod = 30 * time.Second
	if sp := rt.app.Spec.SyncPeriod; sp != nil && sp.Duration > defaultSyncPeriod {
		return sp.Duration
	}
	return defaultSyncPeriod
}

func (rt ReconcileTimer) failureSyncPeriod() time.Duration {
	consecFailures := float64(rt.app.Status.ConsecutiveReconcileFailures)
	// cap consec failures that we are willing to calculate for to avoid overflows
	maxConsecFailures := float64(30)

	if consecFailures > 0 && consecFailures < maxConsecFailures {
		d := time.Duration(math.Exp2(consecFailures)) * time.Second
		if d < rt.syncPeriod() {
			return d
		}
	}

	return rt.syncPeriod()
}

func (rt ReconcileTimer) hasReconcileStatus(c v1alpha1.AppConditionType) bool {
	for _, cond := range rt.app.Status.Conditions {
		if cond.Type == c {
			return true
		}
	}
	return false
}

func (rt ReconcileTimer) applyJitter(t time.Duration) time.Duration {
	const appJitter time.Duration = 5 * time.Second
	return t - appJitter + wait.Jitter(appJitter, 1.0)
}
