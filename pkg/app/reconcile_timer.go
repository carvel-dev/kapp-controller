package app

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

	// If canceled/paused, then no reconcilation until unpaused
	if rt.app.Spec.Canceled || rt.app.Spec.Paused {
		return false
	}

	// Did we deploy at least once?
	lastFetch := rt.app.Status.Fetch
	if lastFetch == nil {
		return true
	}

	if rt.hasReconcileStatus(v1alpha1.ReconcileFailed) {
		if timeAt.UTC().Sub(lastFetch.UpdatedAt.Time) >= rt.failureSyncPeriod() {
			return true
		}
	}

	// Did we deploy too long ago?
	if timeAt.UTC().Sub(lastFetch.UpdatedAt.Time) >= rt.syncPeriod() {
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
	d := time.Duration(math.Exp2(float64(rt.app.Status.ConsecutiveReconcileFailures))) * time.Second
	if d < rt.syncPeriod() {
		return d
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
