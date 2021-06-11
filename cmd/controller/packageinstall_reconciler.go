// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"

	"github.com/go-logr/logr"
	pkgclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/packageinstall"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PackageInstallReconciler struct {
	kcClient  kcclient.Interface
	pkgClient pkgclient.Interface
	log       logr.Logger
}

var _ reconcile.Reconciler = &PackageInstallReconciler{}

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

	return packageinstall.NewPackageInstallCR(existingPkgInstall, log, r.kcClient, r.pkgClient).Reconcile()
}
