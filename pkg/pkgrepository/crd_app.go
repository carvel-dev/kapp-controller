// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"context"
	"fmt"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"reflect"

	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type CRDApp struct {
	app       *App
	appModel  *kcv1alpha1.App
	pkgrModel *pkgingv1alpha1.PackageRepository
	log       logr.Logger
	appClient kcclient.Interface
}

func NewCRDApp(appModel *kcv1alpha1.App, packageRepo *pkgingv1alpha1.PackageRepository, log logr.Logger,
	appClient kcclient.Interface, fetchFactory fetch.Factory,
	templateFactory template.Factory, deployFactory deploy.Factory) *CRDApp {

	crdApp := &CRDApp{appModel: appModel, pkgrModel: packageRepo, log: log, appClient: appClient}

	crdApp.app = NewApp(*appModel, Hooks{
		BlockDeletion:   crdApp.blockDeletion,
		UnblockDeletion: crdApp.unblockDeletion,
		UpdateStatus:    crdApp.updateStatus,
		WatchChanges:    crdApp.watchChanges,
	}, fetchFactory, templateFactory, deployFactory, log)

	return crdApp
}

func (a *CRDApp) blockDeletion() error {
	// Avoid doing unnecessary processing
	if containsString(a.pkgrModel.ObjectMeta.Finalizers, deleteFinalizerName) {
		return nil
	}

	a.log.Info("Blocking deletion")

	return a.updateApp(func(app *pkgingv1alpha1.PackageRepository) {
		if !containsString(app.ObjectMeta.Finalizers, deleteFinalizerName) {
			app.ObjectMeta.Finalizers = append(app.ObjectMeta.Finalizers, deleteFinalizerName)
		}
	})
}

func (a *CRDApp) unblockDeletion() error {
	a.log.Info("Unblocking deletion")
	return a.updateApp(func(app *pkgingv1alpha1.PackageRepository) {
		app.ObjectMeta.Finalizers = removeString(app.ObjectMeta.Finalizers, deleteFinalizerName)
	})
}

func (a *CRDApp) updateStatus(desc string) error {
	a.log.Info("Updating status", "desc", desc)

	var lastErr error
	for i := 0; i < 5; i++ {
		lastErr = a.updateStatusOnce()
		if lastErr == nil {
			return nil
		}
	}

	return lastErr
}

func (a *CRDApp) updateStatusOnce() error {
	existingApp, err := a.appClient.PackagingV1alpha1().PackageRepositories(a.pkgrModel.Namespace).Get(context.Background(), a.pkgrModel.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Fetching app: %s", err)
	}

	if !reflect.DeepEqual(existingApp.Status, a.app.Status()) {

		existingApp.Status = a.app.Status()
		_, err = a.appClient.PackagingV1alpha1().PackageRepositories(existingApp.Namespace).UpdateStatus(context.Background(), existingApp, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *CRDApp) updateApp(updateFunc func(*pkgingv1alpha1.PackageRepository)) error {
	a.log.Info("Updating app")

	existingApp, err := a.appClient.PackagingV1alpha1().PackageRepositories(a.pkgrModel.Namespace).Get(context.Background(), a.pkgrModel.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Updating app: %s", err)
	}

	updateFunc(existingApp)

	_, err = a.appClient.PackagingV1alpha1().PackageRepositories(existingApp.Namespace).Update(context.Background(), existingApp, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Updating app: %s", err)
	}

	return nil
}

func (a *CRDApp) Reconcile(force bool) (reconcile.Result, error) {
	return a.app.Reconcile(force)
}

func (a *CRDApp) watchChanges(callback func(kcv1alpha1.App), cancelCh chan struct{}) error {
	return NewCRDAppWatcher(*a.appModel, a.appClient).Watch(callback, cancelCh)
}

// Get both secret refs/configmap refs
// as single map with all ref entries.
func (a *CRDApp) ResourceRefs() map[reftracker.RefKey]struct{} {
	secrets := a.app.SecretRefs()
	configmaps := a.app.ConfigMapRefs()
	allRefs := map[reftracker.RefKey]struct{}{}

	for secret := range secrets {
		allRefs[secret] = struct{}{}
	}

	for configmap := range configmaps {
		allRefs[configmap] = struct{}{}
	}

	return allRefs
}
