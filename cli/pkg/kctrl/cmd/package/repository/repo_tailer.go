// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdapp "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	kcexternalversions "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/informers/externalversions"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

type RepoTailer struct {
	Namespace string
	Name      string

	ui       ui.UI
	statusUI cmdcore.StatusLoggingUI
	client   kcclient.Interface

	stopperChan chan struct{}
	watchError  error

	lastSeenDeployStdout string
}

func NewRepoTailer(namespace string, name string, ui ui.UI, client kcclient.Interface) *RepoTailer {
	return &RepoTailer{Namespace: namespace, Name: name, ui: ui, statusUI: cmdcore.NewStatusLoggingUI(ui), client: client}
}

func (o *RepoTailer) TailRepoStatus() error {
	o.stopperChan = make(chan struct{})
	_, err := o.client.PackagingV1alpha1().PackageRepositories(o.Namespace).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		if !(errors.IsNotFound(err)) {
			return err
		}
	}

	informerFactory := kcexternalversions.NewFilteredSharedInformerFactory(o.client, 30*time.Minute, o.Namespace, func(opts *metav1.ListOptions) {
		opts.FieldSelector = fmt.Sprintf("metadata.name=%s", o.Name)
	})
	informer := informerFactory.Packaging().V1alpha1().PackageRepositories().Informer()
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
		return fmt.Errorf("Reconciling repository: %s", o.watchError)
	}
	return nil
}

func (o *RepoTailer) stopWatch() {
	close(o.stopperChan)
}

func (o *RepoTailer) printTillCurrent(status kcv1alpha1.AppStatus) error {
	if cmdapp.IsDeleting(status) {
		return nil
	}

	completed, err := cmdapp.NewAppStatusDiff(kcv1alpha1.AppStatus{}, status, o.statusUI).PrintUpdate()
	if err != nil {
		return fmt.Errorf("Reconciling package repository: %s", err)
	}
	if completed {
		o.stopWatch()
	}

	return nil
}

func (o *RepoTailer) udpateEventHandler(oldObj interface{}, newObj interface{}) {
	newRepo, _ := newObj.(*kcpkgv1alpha1.PackageRepository)
	oldRepo, _ := oldObj.(*kcpkgv1alpha1.PackageRepository)

	if newRepo.Generation != newRepo.Status.ObservedGeneration {
		o.statusUI.PrintLogLine(fmt.Sprintf("Waiting for generation %d to be observed", newRepo.Generation), "", false, time.Now())
		return
	}

	mappedOldStatus := o.appStatusFromPkgrStatus(oldRepo.Status)
	mappedNewStatus := o.appStatusFromPkgrStatus(newRepo.Status)

	// o.printUpdate(oldApp.Status, newApp.Status)
	stopWatch, err := cmdapp.NewAppStatusDiff(mappedOldStatus, mappedNewStatus, o.statusUI).PrintUpdate()
	o.watchError = err
	if stopWatch {
		o.stopWatch()
	}
}

func (o *RepoTailer) deleteEventHandler(oldObj interface{}) {
	o.statusUI.PrintLogLine(fmt.Sprintf("Package repository '%s' in namespace '%s' deleted", o.Name, o.Namespace), "", false, time.Now())
	o.stopWatch()
}

func (o *RepoTailer) appStatusFromPkgrStatus(status kcpkgv1alpha1.PackageRepositoryStatus) kcv1alpha1.AppStatus {
	return kcv1alpha1.AppStatus{
		Fetch:    status.Fetch,
		Template: status.Template,
		Deploy:   status.Deploy,
		GenericStatus: kcv1alpha1.GenericStatus{
			Conditions: status.Conditions,
		},
	}
}
