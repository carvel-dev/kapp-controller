// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"math"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type ReconcileTimer struct {
	app v1alpha1.App
}

func NewReconcileTimer(app v1alpha1.App) ReconcileTimer {
	return ReconcileTimer{*app.DeepCopy()}
}

func (rt ReconcileTimer) DurationUntilReady(err error) time.Duration {
	if err != nil || rt.hasReconcileStatus(v1alpha1.ReconcileFailed) || rt.hasReconcileStatus(v1alpha1.DeleteFailed) {
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
	// If canceled/paused, then no reconciliation until unpaused
	if rt.app.Spec.Canceled || rt.app.Spec.Paused {
		return false
	}

	lastReconcileTime, err := rt.lastReconcileTime()
	if err != nil {
		return true
	}

	if rt.hasReconcileStatus(v1alpha1.ReconcileFailed) || rt.hasReconcileStatus(v1alpha1.DeleteFailed) {
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

func (rt ReconcileTimer) hasReconcileStatus(c v1alpha1.ConditionType) bool {
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

func (rt ReconcileTimer) lastReconcileTime() (time.Time, error) {
	// Determine latest time from status and use that as the
	// last reconcile time
	lastReconcileTime := metav1.Time{}
	times := []metav1.Time{}
	if rt.app.Status.Fetch != nil {
		times = append(times, rt.app.Status.Fetch.UpdatedAt)
	}

	if rt.app.Status.Template != nil {
		times = append(times, rt.app.Status.Template.UpdatedAt)
	}

	if rt.app.Status.Deploy != nil {
		times = append(times, rt.app.Status.Deploy.UpdatedAt)
	}

	for _, time := range times {
		if lastReconcileTime.Before(&time) {
			lastReconcileTime = time
		}
	}

	if lastReconcileTime.IsZero() {
		return time.Time{}, fmt.Errorf("could not determine time of last reconcile")
	}

	return lastReconcileTime.Time, nil
}
