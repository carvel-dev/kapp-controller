// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"context"
	"fmt"

	instpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reconciler"

	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PackageRepositoryCR struct {
	model           *instpkgv1alpha1.PackageRepository
	unmodifiedModel *instpkgv1alpha1.PackageRepository

	log    logr.Logger
	client kcclient.Interface
}

func NewPkgRepositoryCR(model *instpkgv1alpha1.PackageRepository, log logr.Logger,
	client kcclient.Interface) *PackageRepositoryCR {

	return &PackageRepositoryCR{model: model, unmodifiedModel: model.DeepCopy(), log: log, client: client}
}

func (ip *PackageRepositoryCR) Reconcile() (reconcile.Result, error) {
	ip.log.Info(fmt.Sprintf("Reconciling PackageRepository '%s'", ip.model.Name))

	if ip.model.DeletionTimestamp != nil {
		// Nothing to do
		return reconcile.Result{}, nil
	}

	// TODO note that we will not be using App CR as a method to download
	// package repositories beyond this quick POC. we would like to decouple storage
	// of packages from kubernetes etcd. Most likely we will rely on k8s API agg layer
	// to serve packages apis directly.

	existingApp, err := ip.client.KappctrlV1alpha1().Apps(appNs).Get(context.Background(), ip.model.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return ip.createApp()
		}
		return reconcile.Result{Requeue: true}, err
	}

	status := &reconciler.Status{
		ip.model.Status.GenericStatus,
		func(st kcv1alpha1.GenericStatus) { ip.model.Status.GenericStatus = st },
	}

	appStatus := reconciler.Status{S: existingApp.Status.GenericStatus}
	ip.log.Info(fmt.Sprintf("Reconciling PackageRepository '%s/%s'", ip.model.Namespace, ip.model.Name))
	switch {
	case appStatus.IsReconciling():
		status.SetReconciling(ip.model.ObjectMeta)
	case appStatus.IsReconcileSucceeded():
		status.SetReconcileCompleted(nil)
	case appStatus.IsReconcileFailed():
		status.SetUsefulErrorMessage(existingApp.Status.UsefulErrorMessage)
		status.SetReconcileCompleted(fmt.Errorf("Syncing packages: (see .status.usefulErrorMessage for details)"))
	}

	result, err := ip.reconcileApp(existingApp)
	if err != nil {
		status.SetReconcileCompleted(err)
	}

	// Always update status
	statusErr := ip.updateStatus()
	if statusErr != nil {
		return reconcile.Result{Requeue: true}, statusErr
	}

	return result, err
}

func (ip *PackageRepositoryCR) createApp() (reconcile.Result, error) {
	desiredApp, err := NewApp(&kcv1alpha1.App{}, ip.model)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	_, err = ip.client.KappctrlV1alpha1().Apps(desiredApp.Namespace).Create(context.Background(), desiredApp, metav1.CreateOptions{})
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	return reconcile.Result{}, nil
}

func (ip *PackageRepositoryCR) reconcileApp(existingApp *kcv1alpha1.App) (reconcile.Result, error) {
	desiredApp, err := NewApp(existingApp, ip.model)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	if !equality.Semantic.DeepEqual(desiredApp, existingApp) {
		_, err = ip.client.KappctrlV1alpha1().Apps(desiredApp.Namespace).Update(context.Background(), desiredApp, metav1.UpdateOptions{})
		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}
	}

	return reconcile.Result{}, nil
}

func (ip *PackageRepositoryCR) updateStatus() error {
	if !equality.Semantic.DeepEqual(ip.unmodifiedModel.Status, ip.model.Status) {
		_, err := ip.client.InstallV1alpha1().PackageRepositories().UpdateStatus(context.Background(), ip.model, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("Updating package repository status: %s", err)
		}
	}
	return nil
}
