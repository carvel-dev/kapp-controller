// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"github.com/go-logr/logr"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/installedpkg"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type InstalledPkgReconciler struct {
	intalledPkgClient kcclient.Interface
	log               logr.Logger
}

var _ reconcile.Reconciler = &InstalledPkgReconciler{}

func (r *InstalledPkgReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	// TODO currently we've decided to get a fresh copy of app so
	// that we do not operate on stale copy for efficiency reasons
	existingInstalledPkg, err := r.intalledPkgClient.InstallV1alpha1().InstalledPackages(request.Namespace).Get(request.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Could not find InstalledPkg", "name", request.Name)
			return reconcile.Result{}, nil // No requeue
		}

		log.Error(err, "Could not fetch InstalledPkg")
		return reconcile.Result{}, err
	}

	return installedpkg.NewInstalledPkgCR(existingInstalledPkg, log, r.intalledPkgClient).Reconcile()
}
