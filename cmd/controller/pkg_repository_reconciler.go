// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"

	"github.com/go-logr/logr"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/pkgrepository"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PkgRepositoryReconciler struct {
	client          kcclient.Interface
	log             logr.Logger
	appFactory      AppFactory
	appRefTracker   *reftracker.AppRefTracker
	appUpdateStatus *reftracker.AppUpdateStatus
}

var _ reconcile.Reconciler = &PkgRepositoryReconciler{}

func NewPkgRepositoryReconciler(appClient kcclient.Interface, log logr.Logger, appFactory AppFactory,
	appRefTracker *reftracker.AppRefTracker, appUpdateStatus *reftracker.AppUpdateStatus) *PkgRepositoryReconciler {
	return &PkgRepositoryReconciler{appClient, log, appFactory, appRefTracker, appUpdateStatus}
}

func (r *PkgRepositoryReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	existingPkgRepository, err := r.client.PackagingV1alpha1().PackageRepositories(request.Namespace).Get(ctx, request.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Could not find PackageRepository", "name", request.Name)
			return reconcile.Result{}, nil // No requeue
		}

		log.Error(err, "Could not fetch PackageRepository")
		return reconcile.Result{}, err
	}

	app, err := pkgrepository.NewPackageRepoApp(existingPkgRepository)
	if err != nil {
		return reconcile.Result{}, err
	}

	crdApp := r.appFactory.NewCRDPackageRepo(app, existingPkgRepository, log)
	r.UpdatePackageRepoRefs(crdApp.ResourceRefs(), app)

	force := false
	pkgrKey := reftracker.NewPackageRepositoryKey(app.Name, app.Namespace)
	if r.appUpdateStatus.IsUpdateNeeded(pkgrKey) {
		r.appUpdateStatus.MarkUpdated(pkgrKey)
		force = true
	}

	return crdApp.Reconcile(force)
}

func (r *PkgRepositoryReconciler) UpdatePackageRepoRefs(refKeys map[reftracker.RefKey]struct{}, app *v1alpha1.App) {
	pkgRepoKey := reftracker.NewPackageRepositoryKey(app.Name, app.Namespace)
	// If PackageRepo is being deleted, remove
	// from all its associated references.
	if app.DeletionTimestamp != nil {
		r.appRefTracker.RemoveAppFromAllRefs(pkgRepoKey)
		return
	}

	// Add new refs for PackageRepo to AppRefTracker/remove
	// any formerly but now unused refs for PackageRepo.
	r.appRefTracker.ReconcileRefs(refKeys, pkgRepoKey)
}
