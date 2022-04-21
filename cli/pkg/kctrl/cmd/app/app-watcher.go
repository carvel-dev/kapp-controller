// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"time"

	"github.com/cppforlife/color"
	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	kcexternalversions "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/informers/externalversions"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

type AppWatcher struct {
	Namespace       string
	Name            string
	Follow          bool
	IgnoreNotExists bool

	ui     ui.UI
	client kcclient.Interface

	lastActiveStage AppStage
}

func NewAppWatcher(namespace string, name string, follow bool, ignoreIfExists bool, ui ui.UI, client kcclient.Interface) *AppWatcher {
	return &AppWatcher{Namespace: namespace, Name: name, Follow: follow, IgnoreNotExists: ignoreIfExists, ui: ui, client: client}
}

func (o *AppWatcher) printTillCurrent(status kcv1alpha1.AppStatus) (AppStage, error) {
	if status.Fetch != nil {
		o.printHeader("Fetch")
		o.printStageMetadata(&status.Fetch.StartedAt, &status.Fetch.UpdatedAt)
		if status.Fetch.ExitCode != 0 && status.Fetch.UpdatedAt.Unix() >= status.Fetch.UpdatedAt.Unix() {
			return fetchStage, fmt.Errorf(status.Fetch.Stderr)
		}
		o.ui.PrintBlock([]byte(status.Fetch.Stdout))
		if status.Fetch.StartedAt.After(status.Fetch.UpdatedAt.Time) {
			o.printOngoing()
			return fetchStage, nil
		}
		o.printFinished()
	}

	if status.Template != nil {
		o.printHeader("Template")
		o.printStageMetadata(nil, &status.Template.UpdatedAt)
		if status.Template.ExitCode != 0 && status.Fetch.StartedAt.Unix() < status.Template.UpdatedAt.Unix() {
			return templateStage, fmt.Errorf(status.Template.Stderr)
		}
		if status.Fetch.StartedAt.After(status.Template.UpdatedAt.Time) {
			o.printOngoing()
			return templateStage, nil
		}
		o.printFinished()
	}

	if status.Deploy != nil {
		o.printHeader("Deploy")
		o.printStageMetadata(&status.Deploy.StartedAt, &status.Deploy.UpdatedAt)
		if status.Deploy.ExitCode != 0 && status.Deploy.StartedAt.Unix() < status.Deploy.UpdatedAt.Unix() {
			return deployStage, fmt.Errorf(status.Deploy.Error)
		}
		o.ui.PrintBlock([]byte(status.Deploy.Stdout))
		if o.hasReconciled(status) {
			o.printFinished()
			return deployStage, nil
		}
		o.printOngoing()
	}

	return reconciled, nil
}

func (o *AppWatcher) PrintTillCurrent(status kcv1alpha1.AppStatus) (AppStage, error) {
	return o.printTillCurrent(status)
}

func (o *AppWatcher) printHeader(header string) {
	o.ui.PrintLinef(color.New(color.Bold).Sprintf("-------------------%s-------------------", header))
}

func (o *AppWatcher) printStageMetadata(startedAt *metav1.Time, updatedAt *metav1.Time) {
	startedAtHeader := uitable.NewHeader("Started At")
	startedAtHeader.Hidden = (startedAt == nil)

	rows := []uitable.Value{
		nil,
		uitable.NewValueTime(updatedAt.Time),
	}

	if startedAt != nil {
		rows = []uitable.Value{
			uitable.NewValueTime(startedAt.Time),
			uitable.NewValueTime(updatedAt.Time),
		}
	}

	table := uitable.Table{
		Transpose: true,

		Header: []uitable.Header{
			startedAtHeader,
			uitable.NewHeader("Updated At"),
		},

		Rows: [][]uitable.Value{rows},
	}

	o.ui.PrintTable(table)
}

func (o *AppWatcher) hasReconciled(status kcv1alpha1.AppStatus) bool {
	for _, condition := range status.Conditions {
		if condition.Type == kcv1alpha1.ReconcileSucceeded && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func (o *AppWatcher) printFinished() {
	o.ui.PrintLinef(color.GreenString("Finished"))
}

func (o *AppWatcher) printOngoing() {
	if !o.Follow {
		o.ui.PrintLinef(color.YellowString("Ongoing"))
	}
}

func (o *AppWatcher) FollowApp(app *kcv1alpha1.App) error {
	currentStage, err := o.printTillCurrent(app.Status)
	if err != nil {
		return err
	}
	o.lastActiveStage = currentStage

	informerFactory := kcexternalversions.NewFilteredSharedInformerFactory(o.client, 30*time.Minute, o.Namespace, func(opts *metav1.ListOptions) {
		opts.FieldSelector = fmt.Sprintf("metadata.name=%s", o.Name)
	})
	informer := informerFactory.Kappctrl().V1alpha1().Apps().Informer()

	stopper := make(chan struct{})
	defer close(stopper)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: o.udpateEventhandler,
	})

	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		return fmt.Errorf("Timed out waiting for caches to sync")
	}
	<-stopper
	return nil
}

func (o *AppWatcher) udpateEventhandler(oldObj interface{}, newObj interface{}) {
	newApp, _ := newObj.(*kcv1alpha1.App)
	// oldApp, _ := oldObj.(*kcv1alpha1.App)

	// if o.reconciliationTriggeredSinceLast(newApp.Status, oldApp.Status) {
	// 	o.printReconciliationDetected()
	// 	currentStage, _ := o.printTillCurrent(newApp.Status)
	// 	o.lastActiveStage = currentStage
	// 	return
	// }

	currentStage, _ := o.printTillCurrent(newApp.Status)
	o.lastActiveStage = currentStage
}

func (o *AppWatcher) printReconciliationDetected() {
	o.ui.PrintLinef(color.New(color.Bold).Sprintf("\nReconciliation triggered at: %s\n", time.Now().String()))
}

func (o *AppWatcher) reconciliationTriggeredSinceLast(oldStatus kcv1alpha1.AppStatus, newStatus kcv1alpha1.AppStatus) bool {
	o.ui.PrintLinef("\n%s\n", oldStatus.Fetch.StartedAt.Before(&newStatus.Fetch.StartedAt))
	if o.lastActiveStage == "" || oldStatus.Fetch.StartedAt.Before(&newStatus.Fetch.StartedAt) {
		return true
	}
	return false
}
