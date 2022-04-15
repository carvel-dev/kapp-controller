// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"fmt"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile is not expected to be called concurrently
func (a *App) Reconcile(force bool) (reconcile.Result, error) {
	defer a.flushUpdateStatus("packagerepository reconciled")

	var err error

	switch {
	case a.app.Spec.Paused:
		a.log.Info("PackageRepository is paused, not reconciling")

		a.markObservedLatest()
		a.app.Status.FriendlyDescription = "Paused"

		err = a.updateStatus("PackageRepository paused")

	case a.app.DeletionTimestamp != nil:
		a.log.Info("Started delete")
		defer func() { a.log.Info("Completed delete") }()

		err = a.reconcileDelete()

	case force || NewReconcileTimer(a.app).IsReadyAt(time.Now()):
		a.log.Info("Started deploy")
		defer func() { a.log.Info("Completed deploy") }()

		err = a.reconcileDeploy()

	default:
		a.log.Info("Reconcile noop")
	}

	return reconcile.Result{RequeueAfter: NewReconcileTimer(a.app).DurationUntilReady(err)}, err
}

func (a *App) reconcileDelete() error {
	a.markObservedLatest()
	a.setDeleting()

	err := a.updateStatus("marking deleting")
	if err != nil {
		return err
	}

	a.resetLastDeployStartedAt()

	result := a.delete()
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

	return a.updateStatus("marking reconcile completed")
}

func (a *App) reconcileFetchTemplateDeploy() exec.CmdRunResult {
	tmpDir := memdir.NewTmpDir("fetch-template-deploy")

	err := tmpDir.Create()
	if err != nil {
		return exec.NewCmdRunResultWithErr(err)
	}

	defer tmpDir.Remove()

	assetsPath := tmpDir.Path()

	{
		a.resetLastFetchStartedAt()

		var fetchResult exec.CmdRunResult
		assetsPath, fetchResult = a.fetch(assetsPath)

		a.app.Status.Fetch = &v1alpha1.AppStatusFetch{
			Stderr:    fetchResult.Stderr,
			Stdout:    fetchResult.Stdout,
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

	tplResult := a.template(assetsPath)

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

	return a.updateLastDeploy(a.deploy(tplResult.Stdout))
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

func (a *App) markObservedLatest() {
	a.app.Status.ObservedGeneration = a.app.Generation
}

func (a *App) setReconciling() {
	a.removeAllConditions()

	a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.Condition{
		Type:   v1alpha1.Reconciling,
		Status: corev1.ConditionTrue,
	})

	a.app.Status.FriendlyDescription = "Reconciling"
}

func (a *App) setReconcileCompleted(result exec.CmdRunResult) {
	a.removeAllConditions()

	if result.Error != nil {
		a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.Condition{
			Type:    v1alpha1.ReconcileFailed,
			Status:  corev1.ConditionTrue,
			Message: result.ErrorStr(),
		})
		a.app.Status.ConsecutiveReconcileFailures++
		a.app.Status.ConsecutiveReconcileSuccesses = 0
		a.app.Status.FriendlyDescription = fmt.Sprintf("Reconcile failed: %s", result.ErrorStr())
		a.setUsefulErrorMessage(result)
	} else {
		a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.Condition{
			Type:    v1alpha1.ReconcileSucceeded,
			Status:  corev1.ConditionTrue,
			Message: "",
		})
		a.app.Status.ConsecutiveReconcileSuccesses++
		a.app.Status.ConsecutiveReconcileFailures = 0
		a.app.Status.FriendlyDescription = "Reconcile succeeded"
		a.app.Status.UsefulErrorMessage = ""
	}
}

func (a *App) setDeleting() {
	a.removeAllConditions()

	a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.Condition{
		Type:   v1alpha1.Deleting,
		Status: corev1.ConditionTrue,
	})

	a.app.Status.FriendlyDescription = "Deleting"
}

func (a *App) setDeleteCompleted(result exec.CmdRunResult) {
	a.removeAllConditions()

	if result.Error != nil {
		a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.Condition{
			Type:    v1alpha1.DeleteFailed,
			Status:  corev1.ConditionTrue,
			Message: result.ErrorStr(),
		})
		a.app.Status.ConsecutiveReconcileFailures++
		a.app.Status.ConsecutiveReconcileSuccesses = 0
		a.app.Status.FriendlyDescription = fmt.Sprintf("Delete failed: %s", result.ErrorStr())
		a.setUsefulErrorMessage(result)
	} else {
		// assume resource will be deleted, hence nothing to update
	}
}

func (a *App) removeAllConditions() {
	a.app.Status.Conditions = nil
}

func (a *App) setUsefulErrorMessage(result exec.CmdRunResult) {
	switch {
	case result.Stderr != "":
		a.app.Status.UsefulErrorMessage = result.Stderr
	default:
		a.app.Status.UsefulErrorMessage = result.ErrorStr()
	}
}
