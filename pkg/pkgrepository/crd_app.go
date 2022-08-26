// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
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
	}, fetchFactory, templateFactory, deployFactory, log)

	return crdApp
}

func (a *CRDApp) blockDeletion() error {
	// Avoid doing unnecessary processing
	if containsString(a.pkgrModel.ObjectMeta.Finalizers, deleteFinalizerName) {
		return nil
	}

	a.log.Info("Blocking deletion")

	return a.updatePackageRepository(func(app *pkgingv1alpha1.PackageRepository) {
		if !containsString(app.ObjectMeta.Finalizers, deleteFinalizerName) {
			app.ObjectMeta.Finalizers = append(app.ObjectMeta.Finalizers, deleteFinalizerName)
		}
	})
}

func (a *CRDApp) unblockDeletion() error {
	a.log.Info("Unblocking deletion")
	return a.updatePackageRepository(func(app *pkgingv1alpha1.PackageRepository) {
		app.ObjectMeta.Finalizers = removeString(app.ObjectMeta.Finalizers, deleteFinalizerName)
		// Need to remove old finalizer that might have been added by previous versions of kapp-controller
		app.ObjectMeta.Finalizers = removeString(app.ObjectMeta.Finalizers, deletePrevFinalizerName)
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
	existingRepo, err := a.appClient.PackagingV1alpha1().PackageRepositories(a.pkgrModel.Namespace).Get(context.Background(), a.pkgrModel.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Fetching PackageRepository: %s", err)
	}

	existingRepo.Status = pkgingv1alpha1.PackageRepositoryStatus{
		Fetch:                         a.app.Status().Fetch,
		Template:                      a.app.Status().Template,
		Deploy:                        a.app.Status().Deploy,
		GenericStatus:                 a.app.Status().GenericStatus,
		ConsecutiveReconcileSuccesses: a.app.Status().ConsecutiveReconcileSuccesses,
		ConsecutiveReconcileFailures:  a.app.Status().ConsecutiveReconcileFailures,
	}

	_, err = a.appClient.PackagingV1alpha1().PackageRepositories(existingRepo.Namespace).UpdateStatus(context.Background(), existingRepo, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (a *CRDApp) updatePackageRepository(updateFunc func(*pkgingv1alpha1.PackageRepository)) error {
	a.log.Info("Updating PackageRepository")

	existingRepo, err := a.appClient.PackagingV1alpha1().PackageRepositories(a.pkgrModel.Namespace).Get(context.Background(), a.pkgrModel.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Updating PackageRepository: %s", err)
	}

	updateFunc(existingRepo)

	_, err = a.appClient.PackagingV1alpha1().PackageRepositories(existingRepo.Namespace).Update(context.Background(), existingRepo, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Updating PackageRepository: %s", err)
	}

	return nil
}

func (a *CRDApp) Reconcile(force bool) (reconcile.Result, error) {
	return a.app.Reconcile(force)
}

func (a *CRDApp) ResourceRefs() map[reftracker.RefKey]struct{} {
	return a.app.SecretRefs()
}
