// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"
	"time"

	"github.com/cppforlife/color"
	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	kcexternalversions "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/informers/externalversions"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

type AppTailer struct {
	Namespace string
	Name      string

	ui       ui.UI
	statusUI cmdcore.StatusLoggingUI
	client   kcclient.Interface

	stopperChan       chan struct{}
	erroredWhileWatch bool
	failureMessage    string
	opts              AppTailerOpts

	lastSeenDeployStdout string
}

type AppTailerOpts struct {
	IgnoreNotExists   bool
	PrintMetadata     bool
	PrintCurrentState bool
}

func NewAppTailer(namespace string, name string, ui ui.UI, client kcclient.Interface, opts AppTailerOpts) *AppTailer {
	return &AppTailer{Namespace: namespace, Name: name, opts: opts, ui: ui, statusUI: cmdcore.NewStatusLoggingUI(ui), client: client}
}

func (o *AppTailer) printTillCurrent(status kcv1alpha1.AppStatus) error {
	if o.isDeleting(status) {
		return nil
	}

	o.printUpdate(kcv1alpha1.AppStatus{}, status)

	if o.erroredWhileWatch {
		return fmt.Errorf("Reconciling app: %s", o.failureMessage)
	}
	return nil
}

func (o *AppTailer) printUpdate(oldStatus kcv1alpha1.AppStatus, status kcv1alpha1.AppStatus) {
	if status.Fetch != nil {
		if oldStatus.Fetch == nil || (!oldStatus.Fetch.StartedAt.Equal(&status.Fetch.StartedAt) && status.Fetch.UpdatedAt.Unix() <= status.Fetch.StartedAt.Unix()) {
			o.statusUI.PrintLogLine("Fetch started", "", false, status.Fetch.StartedAt.Time)
		}
		if oldStatus.Fetch == nil || !oldStatus.Fetch.UpdatedAt.Equal(&status.Fetch.UpdatedAt) {
			if status.Fetch.ExitCode != 0 && status.Fetch.UpdatedAt.Unix() >= status.Fetch.StartedAt.Unix() {
				msg := "Fetch failed"
				o.statusUI.PrintLogLine(msg, status.Fetch.Stderr, true, status.Fetch.UpdatedAt.Time)
				o.failureMessage = msg
				o.stopWatch(true)
				return
			}
			o.statusUI.PrintLogLine("Fetching", status.Fetch.Stdout, false, status.Fetch.UpdatedAt.Time)
			o.statusUI.PrintLogLine("Fetch succeeded", "", false, status.Fetch.UpdatedAt.Time)
		}
	}
	if status.Template != nil {
		if oldStatus.Template == nil || !oldStatus.Template.UpdatedAt.Equal(&status.Template.UpdatedAt) {
			if status.Template.ExitCode != 0 {
				msg := "Template failed"
				o.statusUI.PrintLogLine(msg, status.Template.Stderr, true, status.Template.UpdatedAt.Time)
				o.failureMessage = msg
				o.stopWatch(true)
				return
			}
			o.statusUI.PrintLogLine("Template succeeded", "", false, status.Template.UpdatedAt.Time)
		}
	}
	if status.Deploy != nil {
		isDeleting := o.isDeleting(status)
		ongoingOp := "Deploy"
		if isDeleting {
			ongoingOp = "Delete"
		}
		if oldStatus.Deploy == nil || !oldStatus.Deploy.StartedAt.Equal(&status.Deploy.StartedAt) {
			msg := fmt.Sprintf("%s started", ongoingOp)
			o.statusUI.PrintLogLine(msg, "", false, status.Deploy.StartedAt.Time)
		}
		if oldStatus.Deploy == nil || !oldStatus.Deploy.UpdatedAt.Equal(&status.Deploy.UpdatedAt) {
			if status.Deploy.ExitCode != 0 && status.Deploy.Finished {
				msg := fmt.Sprintf("Deploy failed")
				o.statusUI.PrintLogLine(msg, status.Deploy.Stderr, true, status.Deploy.UpdatedAt.Time)
				o.failureMessage = msg
				o.stopWatch(true)
				return
			}
			o.printDeployStdout(status.Deploy.Stdout, status.Deploy.UpdatedAt.Time, isDeleting)
		}
	}

	if o.hasReconciled(status) {
		o.statusUI.PrintLogLine("App reconciled", "", false, status.Deploy.UpdatedAt.Time)
		o.stopWatch(false)
	}
	failed, errMsg := o.hasFailed(status)
	if failed {
		o.statusUI.PrintLogLine(errMsg, "", true, time.Now())
		o.stopWatch(true)
	}
}

func (o *AppTailer) printInfo(app kcv1alpha1.App) {
	table := uitable.Table{
		Transpose: true,

		Header: []uitable.Header{

			uitable.NewHeader("Name"),
			uitable.NewHeader("Namespace"),
			uitable.NewHeader("Status"),
			uitable.NewHeader("Metrics"),
		},

		Rows: [][]uitable.Value{{
			uitable.NewValueString(app.Name),
			uitable.NewValueString(app.Namespace),
			uitable.NewValueString(o.statusString(app.Status)),
			uitable.NewValueString(o.metricString(app.Status)),
		}},
	}

	o.ui.PrintTable(table)
}

func (o *AppTailer) metricString(status kcv1alpha1.AppStatus) string {
	if status.ConsecutiveReconcileFailures != 0 {
		return fmt.Sprintf("%d consecutive failures", status.ConsecutiveReconcileFailures)
	} else if status.ConsecutiveReconcileSuccesses != 0 {
		return fmt.Sprintf("%d consecutive successes", status.ConsecutiveReconcileSuccesses)
	} else {
		return "0 consecutive failures | 0 consecutive successes"
	}
}

func (o *AppTailer) statusString(status kcv1alpha1.AppStatus) string {
	if len(status.Conditions) < 1 {
		return ""
	}
	for _, condition := range status.Conditions {
		switch condition.Type {
		case kcv1alpha1.ReconcileFailed:
			return color.RedString("Reconcile failed")
		case kcv1alpha1.ReconcileSucceeded:
			return color.GreenString("Reconcile succeeded")
		case kcv1alpha1.DeleteFailed:
			return color.RedString("Deletion failed")
		case kcv1alpha1.Reconciling:
			return "Reconciling"
		case kcv1alpha1.Deleting:
			return "Deleting"
		}
	}
	return status.FriendlyDescription
}

func (o *AppTailer) hasReconciled(status kcv1alpha1.AppStatus) bool {
	for _, condition := range status.Conditions {
		if condition.Type == kcv1alpha1.ReconcileSucceeded && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func (o *AppTailer) hasFailed(status kcv1alpha1.AppStatus) (bool, string) {
	for _, condition := range status.Conditions {
		if condition.Type == kcv1alpha1.ReconcileFailed && condition.Status == corev1.ConditionTrue {
			return true, color.RedString(fmt.Sprintf("%s: %s", kcv1alpha1.ReconcileFailed, status.UsefulErrorMessage))
		}
		if condition.Type == kcv1alpha1.DeleteFailed && condition.Status == corev1.ConditionTrue {
			return true, color.RedString(fmt.Sprintf("%s: %s", kcv1alpha1.DeleteFailed, status.UsefulErrorMessage))
		}
	}
	return false, ""
}

func (o *AppTailer) isDeleting(status kcv1alpha1.AppStatus) bool {
	for _, condition := range status.Conditions {
		if condition.Type == kcv1alpha1.Deleting && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func (o *AppTailer) TailAppStatus() error {
	o.stopperChan = make(chan struct{})
	app, err := o.client.KappctrlV1alpha1().Apps(o.Namespace).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		if !(errors.IsNotFound(err) && o.opts.IgnoreNotExists) {
			return err
		}
	}

	if o.opts.PrintMetadata {
		o.printInfo(*app)
	}

	if o.opts.PrintCurrentState {
		err = o.printTillCurrent(app.Status)
		if err != nil {
			return err
		}

		if o.hasReconciled(app.Status) {
			return nil
		}
	}

	informerFactory := kcexternalversions.NewFilteredSharedInformerFactory(o.client, 30*time.Minute, o.Namespace, func(opts *metav1.ListOptions) {
		opts.FieldSelector = fmt.Sprintf("metadata.name=%s", o.Name)
	})
	informer := informerFactory.Kappctrl().V1alpha1().Apps().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: o.udpateEventHandler,
		DeleteFunc: o.deleteEventHandler,
	})

	go informer.Run(o.stopperChan)
	if !cache.WaitForCacheSync(o.stopperChan, informer.HasSynced) {
		return fmt.Errorf("Timed out waiting for caches to sync")
	}

	<-o.stopperChan
	if o.erroredWhileWatch {
		return fmt.Errorf("Reconciling app: %s", o.failureMessage)
	}
	return nil
}

func (o *AppTailer) stopWatch(failing bool) {
	o.erroredWhileWatch = failing
	close(o.stopperChan)
}

func (o *AppTailer) udpateEventHandler(oldObj interface{}, newObj interface{}) {
	newApp, _ := newObj.(*kcv1alpha1.App)
	oldApp, _ := oldObj.(*kcv1alpha1.App)

	if newApp.Generation != newApp.Status.ObservedGeneration {
		o.statusUI.PrintLogLine(fmt.Sprintf("Waiting for generation %d to be observed", newApp.Generation), "", false, time.Now())
		return
	}

	o.printUpdate(oldApp.Status, newApp.Status)
}

func (o *AppTailer) deleteEventHandler(oldObj interface{}) {
	o.statusUI.PrintLogLine(fmt.Sprintf("App '%s' in namespace '%s' deleted", o.Name, o.Namespace), "", false, time.Now())
	o.stopWatch(false)
}

func (o *AppTailer) printDeployStdout(stdout string, timestamp time.Time, isDeleting bool) {
	if o.lastSeenDeployStdout == "" {
		o.lastSeenDeployStdout = stdout
		msg := "Deploying"
		if isDeleting {
			msg = "Deleting"
		}
		o.statusUI.PrintLogLine(msg, stdout, false, timestamp)
		return
	}

	o.statusUI.PrintMessageBlockDiff(o.lastSeenDeployStdout, stdout, timestamp)

	o.lastSeenDeployStdout = stdout
}
