// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	kappctrlv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	pkgclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type PackageInstallReconciler struct {
	kcClient               kcclient.Interface
	pkgClient              pkgclient.Interface
	coreClient             kubernetes.Interface
	pkgToPkgInstallHandler *PackageInstallVersionHandler
	log                    logr.Logger
}

// NewPackageInstallReconciler is the constructor for the PackageInstallReconciler struct
func NewPackageInstallReconciler(kcClient kcclient.Interface, pkgClient pkgclient.Interface,
	coreClient kubernetes.Interface, pkgToPkgInstallHandler *PackageInstallVersionHandler,
	log logr.Logger) *PackageInstallReconciler {
	return &PackageInstallReconciler{kcClient, pkgClient, coreClient, pkgToPkgInstallHandler, log}
}

var _ reconcile.Reconciler = &PackageInstallReconciler{}

func (r *PackageInstallReconciler) AttachWatches(controller controller.Controller) error {
	err := controller.Watch(&source.Kind{Type: &pkgingv1alpha1.PackageInstall{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return fmt.Errorf("Watching PackageInstalls: %s", err)
	}

	err = controller.Watch(&source.Kind{Type: &datapkgingv1alpha1.Package{}}, r.pkgToPkgInstallHandler)
	if err != nil {
		return fmt.Errorf("Watching Packages: %s", err)
	}

	err = controller.Watch(&source.Kind{Type: &kappctrlv1alpha1.App{}}, &handler.EnqueueRequestForOwner{
		OwnerType:    &pkgingv1alpha1.PackageInstall{},
		IsController: true,
	})
	if err != nil {
		return fmt.Errorf("Watching Apps: %s", err)
	}

	return nil
}

func (r *PackageInstallReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	existingPkgInstall, err := r.kcClient.PackagingV1alpha1().PackageInstalls(request.Namespace).Get(ctx, request.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Could not find PackageInstall", "name", request.Name)
			return reconcile.Result{}, nil // No requeue
		}

		log.Error(err, "Could not fetch PackageInstall")
		return reconcile.Result{}, err
	}

	return NewPackageInstallCR(existingPkgInstall, log, r.kcClient, r.pkgClient, r.coreClient).Reconcile()
}
