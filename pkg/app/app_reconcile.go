package app

import (
	"time"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/exec"
	"github.com/k14s/kapp-controller/pkg/memdir"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *App) Reconcile() error {
	isDeleting := a.app.DeletionTimestamp != nil

	if isDeleting || a.shouldReconcile() {
		var result exec.CmdRunResult

		a.setReconciling()
		a.updateStatus()

		defer func() {
			a.setReconcileCompleted(result)
			a.updateStatus()
		}()

		// TODO a.app.Status.ManagedAppName = a.managedName()

		if isDeleting {
			a.resetLastDeployStartedAt()
			result = a.delete(a.updateLastDeployNoReturn)
			return result.Error
		}

		if !a.app.Spec.Paused {
			result = a.reconcileFetchTemplateDeploy()
			// Reconcile inspect regardless of deploy success
			a.reconcileInspect()
		}
	}

	return nil
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

		a.updateStatus()

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

	a.updateStatus()

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

	a.updateStatus()

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

	a.updateStatus()

	return nil
}

func (a *App) setReconciling() {
	a.removeCondition(v1alpha1.Reconciling)
	a.removeCondition(v1alpha1.ReconcileFailed)

	a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.AppCondition{
		Type:   v1alpha1.Reconciling,
		Status: corev1.ConditionTrue,
	})

	a.app.Status.ObservedGeneration = a.app.Generation
}

func (a *App) setReconcileCompleted(result exec.CmdRunResult) {
	a.removeCondition(v1alpha1.Reconciling)
	a.removeCondition(v1alpha1.ReconcileFailed)

	if result.Error != nil {
		a.app.Status.Conditions = append(a.app.Status.Conditions, v1alpha1.AppCondition{
			Type:    v1alpha1.ReconcileFailed,
			Status:  corev1.ConditionTrue,
			Message: result.ErrorStr(),
		})
	}
}

func (a *App) shouldReconcile() bool {
	// Did resource spec change?
	if a.app.Status.ObservedGeneration != a.app.Generation {
		return true
	}

	// Did previous deploy fail?
	for _, cond := range a.app.Status.Conditions {
		if cond.Type == v1alpha1.ReconcileFailed {
			return true
		}
	}

	// Did we deploy at least once?
	lastDeploy := a.app.Status.Deploy
	if lastDeploy == nil {
		return true
	}

	// Did we deploy some time ago?
	if time.Now().UTC().Sub(lastDeploy.UpdatedAt.Time) > 30*time.Second {
		return true
	}

	return false
}

func (a *App) removeCondition(type_ v1alpha1.AppConditionType) {
	for i, cond := range a.app.Status.Conditions {
		if cond.Type == type_ {
			a.app.Status.Conditions = append(a.app.Status.Conditions[:i], a.app.Status.Conditions[i+1:]...)
			return
		}
	}
}
