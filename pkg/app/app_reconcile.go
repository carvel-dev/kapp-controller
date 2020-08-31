// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"time"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/exec"
	"github.com/k14s/kapp-controller/pkg/memdir"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile is not expected to be called concurrently
func (a *App) Reconcile() (reconcile.Result, error) {
	defer a.flushUpdateStatus("app reconciled")

	var err error

	switch {
	case a.app.Spec.Canceled || a.app.Spec.Paused:
		a.log.Info("App is canceled or paused, not reconciling")

		a.markObservedLatest()
		a.app.Status.FriendlyDescription = "Canceled/paused"

		err = a.updateStatus("app canceled/paused")

	case a.app.DeletionTimestamp != nil:
		a.log.Info("Started delete")
		defer func() { a.log.Info("Completed delete") }()

		err = a.reconcileDelete()

	case a.shouldReconcile(time.Now()):
		a.log.Info("Started deploy")
		defer func() { a.log.Info("Completed deploy") }()

		err = a.reconcileDeploy()

	default:
		a.log.Info("Reconcile noop")
	}

	return a.requeueIfNecessary(), err
}

func (a *App) reconcileDelete() error {
	a.markObservedLatest()
	a.setDeleting()

	err := a.updateStatus("marking deleting")
	if err != nil {
		return err
	}

	a.resetLastDeployStartedAt()

	result := a.delete(a.updateLastDeployNoReturn)
	a.setDeleteCompleted(result)

	// Resource is gone so this will error, ignore it
	_ = a.updateStatus("marking delete completed")
	return nil
}

func (a *App) reconcileDeploy() error {
	a.markObservedLatest()
	a.setReconciling()

	err := a.updateStatus("marking reconciling")
	if err != nil {
		return err
	}

	result := a.reconcileFetchTemplateDeploy()
	a.setReconcileCompleted(result)

	// Reconcile inspect regardless of deploy success
	_ = a.reconcileInspect()

	return a.updateStatus("marking reconcile completed")
}

func (a *App) reconcileFetchTemplateDeploy() exec.CmdRunResult {
	tmpDir := memdir.NewTmpDir("fetch-template-deploy")

	err := tmpDir.Create()
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	defer tmpDir.Remove()

	{
		a.resetLastFetchStartedAt()

		fetchResult := a.fetch(tmpDir.Path())

		a.app.Status.Fetch = &v1alpha1.AppStatusFetch{
			Stderr:    fetchResult.Stderr,
			ExitCode:  fetchResult.ExitCode,
			Error:     fetchResult.ErrorStr(),
			StartedAt: a.app.Status.Fetch.StartedAt,
			UpdatedAt: metav1.NewTime(time.Now().UTC()),
		}

		err := a.updateStatus("marking fetch completed")
		if err != nil {
			return exec.NewCmdRunResultWithErr(err)
		}

		if fetchResult.Error != nil {
			return fetchResult
		}
	}

	tplResult := a.template(tmpDir.Path())

	a.app.Status.Template = &v1alpha1.AppStatusTemplate{
		Stderr:    tplResult.Stderr,
		ExitCode:  tplResult.ExitCode,
		Error:     tplResult.ErrorStr(),
		UpdatedAt: metav1.NewTime(time.Now().UTC()),
	}

	err = a.updateStatus("marking template completed")
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	if tplResult.Error != nil {
		return tplResult
	}

	a.resetLastDeployStartedAt()

	return a.updateLastDeploy(a.deploy(tplResult.Stdout, a.updateLastDeployNoReturn))
}

func (a *App) updateLastDeploy(result exec.CmdRunResult) exec.CmdRunResult {
	result = result.WithFriendlyYAMLStrings()

	a.app.Status.Deploy = &v1alpha1.AppStatusDeploy{
		Stdout:    result.Stdout,
		Stderr:    result.Stderr,
		Finished:  result.Finished,
		ExitCode:  result.ExitCode,
		Error:     result.ErrorStr(),
		StartedAt: a.app.Status.Deploy.StartedAt,
		UpdatedAt: metav1.NewTime(time.Now().UTC()),
	}

	a.updateStatus("marking last deploy")

	return result
}

func (a *App) updateLastDeployNoReturn(result exec.CmdRunResult) {
	a.updateLastDeploy(result)
}

func (a *App) resetLastFetchStartedAt() {
	if a.app.Status.Fetch == nil {
		a.app.Status.Fetch = &v1alpha1.AppStatusFetch{}
	}
	a.app.Status.Fetch.StartedAt = metav1.NewTime(time.Now().UTC())
}

func (a *App) resetLastDeployStartedAt() {
	if a.app.Status.Deploy == nil {
		a.app.Status.Deploy = &v1alpha1.AppStatusDeploy{}
	}
	a.app.Status.Deploy.StartedAt = metav1.NewTime(time.Now().UTC())
}

func (a *App) reconcileInspect() error {
	inspectResult := a.inspect().WithFriendlyYAMLStrings()

	a.app.Status.Inspect = &v1alpha1.AppStatusInspect{
		Stdout:    inspectResult.Stdout,
		Stderr:    inspectResult.Stderr,
		ExitCode:  inspectResult.ExitCode,
		Error:     inspectResult.ErrorStr(),
		UpdatedAt: metav1.NewTime(time.Now().UTC()),
	}

	return a.updateStatus("marking inspect completed")
}

func (a *App) markObservedLatest() {
	a.app.Status.ObservedGeneration = a.app.Generation
}

func (a *App) setReconciling() {
	a.removeAllConditions()

	a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.AppCondition{
		Type:   v1alpha1.Reconciling,
		Status: corev1.ConditionTrue,
	})

	a.app.Status.FriendlyDescription = "Reconciling"
}

func (a *App) setReconcileCompleted(result exec.CmdRunResult) {
	a.removeAllConditions()

	if result.Error != nil {
		a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.AppCondition{
			Type:    v1alpha1.ReconcileFailed,
			Status:  corev1.ConditionTrue,
			Message: result.ErrorStr(),
		})
		a.app.Status.FriendlyDescription = fmt.Sprintf("Reconcile failed: %s", result.ErrorStr())
	} else {
		a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.AppCondition{
			Type:    v1alpha1.ReconcileSucceeded,
			Status:  corev1.ConditionTrue,
			Message: "",
		})
		a.app.Status.FriendlyDescription = "Reconcile succeeded"
	}
}

func (a *App) setDeleting() {
	a.removeAllConditions()

	a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.AppCondition{
		Type:   v1alpha1.Deleting,
		Status: corev1.ConditionTrue,
	})

	a.app.Status.FriendlyDescription = "Deleting"
}

func (a *App) setDeleteCompleted(result exec.CmdRunResult) {
	a.removeAllConditions()

	if result.Error != nil {
		a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.AppCondition{
			Type:    v1alpha1.DeleteFailed,
			Status:  corev1.ConditionTrue,
			Message: result.ErrorStr(),
		})
		a.app.Status.FriendlyDescription = fmt.Sprintf("Delete failed: %s", result.ErrorStr())
	} else {
		// assume resource will be deleted, hence nothing to update
	}
}

func (a *App) syncPeriod() time.Duration {
	const DefaultSyncPeriod = 30 * time.Second
	if sp := a.app.Spec.SyncPeriod; sp != nil && sp.Duration > DefaultSyncPeriod {
		return sp.Duration
	}
	return DefaultSyncPeriod
}

func (a *App) requeueIfNecessary() reconcile.Result {
	var (
		shortDelay = 4 * time.Second
		// Must always be >= tooLongAfterSuccess so that we dont requeue
		// without work to do
		// replace last 5 seconds with int from range [5,10]
		longerDelay = a.syncPeriod() - 5 + wait.Jitter(5*time.Second, 1.0)
	)

	if a.shouldReconcile(time.Now().Add(shortDelay)) {
		return reconcile.Result{RequeueAfter: shortDelay}
	}
	return reconcile.Result{RequeueAfter: longerDelay}
}

func (a *App) shouldReconcile(timeAt time.Time) bool {
	const (
		tooLongAfterFailure = 3 * time.Second
	)
	tooLongAfterSuccess := a.syncPeriod()

	// Did resource spec change?
	if a.app.Status.ObservedGeneration != a.app.Generation {
		return true
	}

	// If canceled/paused, then no reconcilation until unpaused
	if a.app.Spec.Canceled || a.app.Spec.Paused {
		return false
	}

	// Did we deploy at least once?
	lastDeploy := a.app.Status.Deploy
	if lastDeploy == nil {
		return true
	}

	// Did previous deploy fail?
	for _, cond := range a.app.Status.Conditions {
		if cond.Type == v1alpha1.ReconcileFailed {
			// Did we try too long ago?
			if timeAt.UTC().Sub(lastDeploy.UpdatedAt.Time) > tooLongAfterFailure {
				return true
			}
		}
	}

	// Did we deploy too long ago?
	if timeAt.UTC().Sub(lastDeploy.UpdatedAt.Time) > tooLongAfterSuccess {
		return true
	}

	return false
}

func (a *App) removeAllConditions() {
	a.app.Status.Conditions = nil
}

func (a *App) removeCondition(type_ v1alpha1.AppConditionType) {
	for i, cond := range a.app.Status.Conditions {
		if cond.Type == type_ {
			a.app.Status.Conditions = append(a.app.Status.Conditions[:i], a.app.Status.Conditions[i+1:]...)
			return
		}
	}
}
