// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"testing"
	"time"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	defaultTimerOpts = ReconcileTimerOpts{
		DefaultSyncPeriod: 30 * time.Second,
		MinimumSyncPeriod: 30 * time.Second,
	}
)

func TestSucceededDurationUntilReady(t *testing.T) {
	syncPeriod := 1 * time.Minute
	app := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			SyncPeriod: &metav1.Duration{Duration: syncPeriod},
		},
		Status: v1alpha1.AppStatus{
			GenericStatus: v1alpha1.GenericStatus{
				Conditions: []v1alpha1.Condition{{Type: v1alpha1.ReconcileSucceeded}},
			},
		},
	}

	for i := 0; i < 100; i++ {
		durationUntilReady := NewReconcileTimer(app, defaultTimerOpts).DurationUntilReady(nil)
		assert.False(
			t,
			durationUntilReady < syncPeriod || durationUntilReady > (syncPeriod+10*time.Second),
			"Expected duration until next reconcile to be in [syncPeriod, syncPeriod + 10]",
		)
	}
}

func TestFailureSyncMathOverflowGuard(t *testing.T) {
	syncPeriod := 30 * time.Second
	app := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			SyncPeriod: &metav1.Duration{Duration: syncPeriod},
		},
		Status: v1alpha1.AppStatus{
			ConsecutiveReconcileFailures: 2700, // number so large 2^x will definitely overflow
			GenericStatus: v1alpha1.GenericStatus{
				Conditions: []v1alpha1.Condition{{Type: v1alpha1.ReconcileFailed}},
			},
		},
	}

	delay := NewReconcileTimer(app, defaultTimerOpts).DurationUntilReady(nil)

	// In the overflow case, delay would be negative due to the 2^x overflow,
	// which would be less than syncPeriod and would be returned. This checks
	// the guard against overflow works
	assert.Equal(t, delay, syncPeriod, "Expected failureSync period to handle an overflow")
}

func TestConsecutiveFailuresOverflowGuard(t *testing.T) {
	syncPeriod := 30 * time.Second
	app := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			SyncPeriod: &metav1.Duration{Duration: syncPeriod},
		},
		Status: v1alpha1.AppStatus{
			ConsecutiveReconcileFailures: -2, // number so large 2^x will definitely overflow
			GenericStatus: v1alpha1.GenericStatus{
				Conditions: []v1alpha1.Condition{{Type: v1alpha1.ReconcileFailed}},
			},
		},
	}

	delay := NewReconcileTimer(app, defaultTimerOpts).DurationUntilReady(nil)

	// make sure that if consecutive failed reconciles has overflowed, we just
	// return syncPeriod instead of a fractional duration (due to neg exp in 2^x)
	assert.Equal(t, delay, syncPeriod, "Expected failureSync period to handle an overflow")
}

func TestFailedDurationUntilReady(t *testing.T) {
	syncPeriod := 30 * time.Second
	app := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			SyncPeriod: &metav1.Duration{Duration: syncPeriod},
		},
		Status: v1alpha1.AppStatus{
			GenericStatus: v1alpha1.GenericStatus{
				Conditions: []v1alpha1.Condition{{Type: v1alpha1.ReconcileFailed}},
			},
		},
	}

	type measurement struct {
		NumberOfFailedReconciles int
		ExpectedDuration         time.Duration
	}

	measurements := []measurement{
		{NumberOfFailedReconciles: 1, ExpectedDuration: 2 * time.Second},
		{NumberOfFailedReconciles: 2, ExpectedDuration: 4 * time.Second},
		{NumberOfFailedReconciles: 3, ExpectedDuration: 8 * time.Second},
		{NumberOfFailedReconciles: 4, ExpectedDuration: 16 * time.Second},
		{NumberOfFailedReconciles: 5, ExpectedDuration: 30 * time.Second},
		{NumberOfFailedReconciles: 6, ExpectedDuration: 30 * time.Second},
	}

	for _, m := range measurements {
		app.Status.ConsecutiveReconcileFailures = m.NumberOfFailedReconciles

		durationUntilReady := NewReconcileTimer(app, defaultTimerOpts).DurationUntilReady(nil)
		assert.Equalf(
			t,
			m.ExpectedDuration,
			durationUntilReady,
			"Expected app with %d failure(s) to have duration %d but got %d",
			m.NumberOfFailedReconciles,
			m.ExpectedDuration,
			durationUntilReady,
		)
	}
}

func TestSucceededIsReadyAt(t *testing.T) {
	syncPeriod := 30 * time.Second
	timeNow := time.Now()
	timeOfReady := timeNow.Add(syncPeriod)

	app := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			SyncPeriod: &metav1.Duration{Duration: syncPeriod},
		},
		Status: v1alpha1.AppStatus{
			Fetch: &v1alpha1.AppStatusFetch{
				UpdatedAt: metav1.Time{Time: timeNow},
			},
			GenericStatus: v1alpha1.GenericStatus{
				Conditions: []v1alpha1.Condition{{Type: v1alpha1.ReconcileSucceeded}},
			},
		},
	}

	isReady := NewReconcileTimer(app, defaultTimerOpts).IsReadyAt(timeOfReady)
	assert.True(t, isReady, "Expected app to be ready after syncPeriod of 30s")

	isReady = NewReconcileTimer(app, defaultTimerOpts).IsReadyAt(timeOfReady.Add(1 * time.Second))
	assert.True(t, isReady, "Expected app to be ready after exceeding syncPeriod of 30s")

	isReady = NewReconcileTimer(app, defaultTimerOpts).IsReadyAt(timeOfReady.Add(-1 * time.Second))
	assert.False(t, isReady, "Expected app to not be ready under syncPeriod of 30s")
}

func TestFailedIsReadyAt(t *testing.T) {
	syncPeriod := 2 * time.Second
	timeNow := time.Now()
	timeOfReady := timeNow.Add(syncPeriod)

	app := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			SyncPeriod: &metav1.Duration{Duration: syncPeriod},
		},
		Status: v1alpha1.AppStatus{
			Fetch: &v1alpha1.AppStatusFetch{
				UpdatedAt: metav1.Time{Time: timeNow},
			},
			GenericStatus: v1alpha1.GenericStatus{
				Conditions: []v1alpha1.Condition{{Type: v1alpha1.ReconcileFailed}},
			},
			ConsecutiveReconcileFailures: 1,
		},
	}

	isReady := NewReconcileTimer(app, defaultTimerOpts).IsReadyAt(timeOfReady)
	assert.True(t, isReady, "Expected app to be ready after syncPeriod of 2s")

	isReady = NewReconcileTimer(app, defaultTimerOpts).IsReadyAt(timeOfReady.Add(1 * time.Second))
	assert.True(t, isReady, "Expected app to be ready after exceeding syncPeriod of 2s")

	isReady = NewReconcileTimer(app, defaultTimerOpts).IsReadyAt(timeOfReady.Add(-1 * time.Second))
	assert.False(t, isReady, "Expected app to not be ready under syncPeriod of 2s")
}

func TestIsReadyAtWithStaleDeployTime(t *testing.T) {
	syncPeriod := 2 * time.Second
	timeNow := time.Now()
	timeOfReady := timeNow.Add(syncPeriod)

	app := v1alpha1.App{
		Spec: v1alpha1.AppSpec{
			SyncPeriod: &metav1.Duration{Duration: syncPeriod},
		},
		Status: v1alpha1.AppStatus{
			Fetch: &v1alpha1.AppStatusFetch{
				UpdatedAt: metav1.Time{Time: timeOfReady},
				Error:     "I've failed you",
			},
			Deploy: &v1alpha1.AppStatusDeploy{
				UpdatedAt: metav1.Time{Time: timeNow},
			},
			ConsecutiveReconcileFailures: 1,
			GenericStatus: v1alpha1.GenericStatus{
				Conditions: []v1alpha1.Condition{{Type: v1alpha1.ReconcileFailed}},
			},
		},
	}

	isReady := NewReconcileTimer(app, defaultTimerOpts).IsReadyAt(timeOfReady.Add(1 * time.Second))
	require.False(t, isReady, "Expected app not to be ready, because deploy time is stale")
}
