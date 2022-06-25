// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	kcexternalversions "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/informers/externalversions"
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

	stopperChan chan struct{}
	watchError  error
	opts        AppTailerOpts

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
	if IsDeleting(status) {
		return nil
	}

	completed, err := NewAppStatusDiff(kcv1alpha1.AppStatus{}, status, o.statusUI).PrintUpdate()
	if err != nil {
		return fmt.Errorf("Reconciling app: %s", err)
	}
	if completed {
		o.stopWatch()
	}

	return nil
}

func (o *AppTailer) printInfo(app kcv1alpha1.App) {
	status, isFailing := appStatus(&app)
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
			uitable.ValueFmt{V: uitable.NewValueString(status), Error: isFailing},
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

		if HasReconciled(app.Status) {
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
	if o.watchError != nil {
		return fmt.Errorf("Reconciling app: %s", o.watchError)
	}
	return nil
}

func (o *AppTailer) stopWatch() {
	close(o.stopperChan)
}

func (o *AppTailer) udpateEventHandler(oldObj interface{}, newObj interface{}) {
	newApp, _ := newObj.(*kcv1alpha1.App)
	oldApp, _ := oldObj.(*kcv1alpha1.App)

	if newApp.Generation != newApp.Status.ObservedGeneration {
		o.statusUI.PrintLogLine(fmt.Sprintf("Waiting for generation %d to be observed", newApp.Generation), "", false, time.Now())
		return
	}

	// o.printUpdate(oldApp.Status, newApp.Status)
	stopWatch, err := NewAppStatusDiff(oldApp.Status, newApp.Status, o.statusUI).PrintUpdate()
	o.watchError = err
	if stopWatch {
		o.stopWatch()
	}
}

func (o *AppTailer) deleteEventHandler(oldObj interface{}) {
	o.statusUI.PrintLogLine(fmt.Sprintf("App '%s' in namespace '%s' deleted", o.Name, o.Namespace), "", false, time.Now())
	o.stopWatch()
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

type AppStatusDiff struct {
	old kcv1alpha1.AppStatus
	new kcv1alpha1.AppStatus

	statusUI cmdcore.StatusLoggingUI

	lastSeenDeployStdout string
}

func NewAppStatusDiff(old kcv1alpha1.AppStatus, new kcv1alpha1.AppStatus, statusUI cmdcore.StatusLoggingUI) *AppStatusDiff {
	return &AppStatusDiff{old: old, new: new, statusUI: statusUI}
}

func (d *AppStatusDiff) PrintUpdate() (bool, error) {
	if d.new.Fetch != nil {
		if d.old.Fetch == nil || (!d.old.Fetch.StartedAt.Equal(&d.new.Fetch.StartedAt) && d.new.Fetch.UpdatedAt.Unix() <= d.new.Fetch.StartedAt.Unix()) {
			d.statusUI.PrintLogLine("Fetch started", "", false, d.new.Fetch.StartedAt.Time)
		}
		if d.old.Fetch == nil || !d.old.Fetch.UpdatedAt.Equal(&d.new.Fetch.UpdatedAt) {
			if d.new.Fetch.ExitCode != 0 && d.new.Fetch.UpdatedAt.Unix() >= d.new.Fetch.StartedAt.Unix() {
				msg := "Fetch failed"
				errLog := d.new.Fetch.Stderr + "\n" + d.new.Fetch.Error
				d.statusUI.PrintLogLine(msg, errLog, true, d.new.Fetch.UpdatedAt.Time)
				return true, fmt.Errorf(msg)
			}
			d.statusUI.PrintLogLine("Fetching", d.new.Fetch.Stdout, false, d.new.Fetch.UpdatedAt.Time)
			d.statusUI.PrintLogLine("Fetch succeeded", "", false, d.new.Fetch.UpdatedAt.Time)
		}
	}
	if d.new.Template != nil {
		if d.old.Template == nil || !d.old.Template.UpdatedAt.Equal(&d.new.Template.UpdatedAt) {
			if d.new.Template.ExitCode != 0 {
				msg := "Template failed"
				errLog := d.new.Template.Stderr + "\n" + d.new.Template.Error
				d.statusUI.PrintLogLine(msg, errLog, true, d.new.Template.UpdatedAt.Time)
				return true, fmt.Errorf(msg)
			}
			d.statusUI.PrintLogLine("Template succeeded", "", false, d.new.Template.UpdatedAt.Time)
		}
	}
	if d.new.Deploy != nil {
		isDeleting := IsDeleting(d.new)
		ongoingOp := "Deploy"
		if isDeleting {
			ongoingOp = "Delete"
		}
		if d.old.Deploy == nil || !d.old.Deploy.StartedAt.Equal(&d.new.Deploy.StartedAt) {
			msg := fmt.Sprintf("%s started", ongoingOp)
			d.statusUI.PrintLogLine(msg, "", false, d.new.Deploy.StartedAt.Time)
		}
		if d.old.Deploy == nil || !d.old.Deploy.UpdatedAt.Equal(&d.new.Deploy.UpdatedAt) {
			if d.new.Deploy.ExitCode != 0 && d.new.Deploy.Finished {
				msg := fmt.Sprintf("%s failed", ongoingOp)
				errLog := d.new.Deploy.Stderr + "\n" + d.new.Deploy.Error
				d.statusUI.PrintLogLine(msg, errLog, true, d.new.Deploy.UpdatedAt.Time)
				return true, fmt.Errorf(msg)
			}
			d.printDeployStdout(d.new.Deploy.Stdout, d.new.Deploy.UpdatedAt.Time, isDeleting)
		}
	}

	if HasReconciled(d.new) {
		d.statusUI.PrintLogLine("Deploy succeeded", "", false, d.new.Deploy.UpdatedAt.Time)
		return true, nil
	}
	failed, errMsg := HasFailed(d.new)
	if failed {
		d.statusUI.PrintLogLine(errMsg, "", true, time.Now())
		return true, fmt.Errorf(errMsg)
	}
	return false, nil
}

func (d *AppStatusDiff) printDeployStdout(stdout string, timestamp time.Time, isDeleting bool) {
	if d.lastSeenDeployStdout == "" {
		d.lastSeenDeployStdout = stdout
		msg := "Deploying"
		if isDeleting {
			msg = "Deleting"
		}
		d.statusUI.PrintLogLine(msg, stdout, false, timestamp)
		return
	}

	d.statusUI.PrintMessageBlockDiff(d.lastSeenDeployStdout, stdout, timestamp)

	d.lastSeenDeployStdout = stdout
}
