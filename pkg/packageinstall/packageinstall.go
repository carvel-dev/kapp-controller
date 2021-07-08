// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	pkgclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reconciler"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// nolint: revive
type PackageInstallCR struct {
	model           *pkgingv1alpha1.PackageInstall
	unmodifiedModel *pkgingv1alpha1.PackageInstall

	log       logr.Logger
	kcclient  kcclient.Interface
	pkgclient pkgclient.Interface
}

func NewPackageInstallCR(model *pkgingv1alpha1.PackageInstall, log logr.Logger,
	kcclient kcclient.Interface, pkgclient pkgclient.Interface) *PackageInstallCR {

	return &PackageInstallCR{model: model, unmodifiedModel: model.DeepCopy(), log: log, kcclient: kcclient, pkgclient: pkgclient}
}

func (pi *PackageInstallCR) Reconcile() (reconcile.Result, error) {
	pi.log.Info(fmt.Sprintf("Reconciling PackageInstall '%s/%s'", pi.model.Namespace, pi.model.Name))

	status := &reconciler.Status{
		pi.model.Status.GenericStatus,
		func(st kcv1alpha1.GenericStatus) { pi.model.Status.GenericStatus = st },
	}

	var result reconcile.Result
	var err error

	if pi.model.DeletionTimestamp != nil {
		result, err = pi.reconcileDelete(status)
		if err != nil {
			status.SetDeleteCompleted(err)
		}
	} else {
		result, err = pi.reconcile(status)
		if err != nil {
			status.SetReconcileCompleted(err)
		}
	}

	// Always update status
	statusErr := pi.updateStatus()
	if statusErr != nil {
		return reconcile.Result{Requeue: true}, statusErr
	}

	return result, err
}

func (pi *PackageInstallCR) reconcile(modelStatus *reconciler.Status) (reconcile.Result, error) {
	err := pi.blockDeletion()
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	modelStatus.SetReconciling(pi.model.ObjectMeta)

	pv, err := pi.referencedPkgVersion()
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	pi.model.Status.Version = pv.Spec.Version

	existingApp, err := pi.kcclient.KappctrlV1alpha1().Apps(pi.model.Namespace).Get(context.Background(), pi.model.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return pi.createAppFromPackage(pv)
		}
		return reconcile.Result{Requeue: true}, err
	}

	appStatus := reconciler.Status{S: existingApp.Status.GenericStatus}
	pi.log.Info(fmt.Sprintf("Reconciling PackageInstall '%s/%s'", pi.model.Namespace, pi.model.Name))
	switch {
	case appStatus.IsReconciling():
		modelStatus.SetReconciling(pi.model.ObjectMeta)
	case appStatus.IsReconcileSucceeded():
		modelStatus.SetReconcileCompleted(nil)
	case appStatus.IsReconcileFailed():
		modelStatus.SetUsefulErrorMessage(existingApp.Status.UsefulErrorMessage)
		modelStatus.SetReconcileCompleted(fmt.Errorf("Error (see .status.usefulErrorMessage for details)"))
	}

	return pi.reconcileAppWithPackage(existingApp, pv)
}

func (pi *PackageInstallCR) createAppFromPackage(pv datapkgingv1alpha1.Package) (reconcile.Result, error) {
	desiredApp, err := NewApp(&v1alpha1.App{}, pi.model, pv)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	_, err = pi.kcclient.KappctrlV1alpha1().Apps(desiredApp.Namespace).Create(context.Background(), desiredApp, metav1.CreateOptions{})
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	return reconcile.Result{}, nil
}

func (pi *PackageInstallCR) reconcileAppWithPackage(existingApp *kcv1alpha1.App, pv datapkgingv1alpha1.Package) (reconcile.Result, error) {
	desiredApp, err := NewApp(existingApp, pi.model, pv)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	if !equality.Semantic.DeepEqual(desiredApp, existingApp) {
		_, err = pi.kcclient.KappctrlV1alpha1().Apps(desiredApp.Namespace).Update(context.Background(), desiredApp, metav1.UpdateOptions{})
		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}
	}

	return reconcile.Result{}, nil
}

func (pi *PackageInstallCR) referencedPkgVersion() (datapkgingv1alpha1.Package, error) {
	if pi.model.Spec.PackageRef == nil {
		return datapkgingv1alpha1.Package{}, fmt.Errorf("Expected non nil PackageRef")
	}

	semverConfig := pi.model.Spec.PackageRef.VersionSelection

	pvList, err := pi.pkgclient.DataV1alpha1().Packages(pi.model.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return datapkgingv1alpha1.Package{}, err
	}

	var versionStrs []string
	versionToPkg := map[string]datapkgingv1alpha1.Package{}

	for _, pv := range pvList.Items {
		if pv.Spec.RefName == pi.model.Spec.PackageRef.RefName {
			versionStrs = append(versionStrs, pv.Spec.Version)
			versionToPkg[pv.Spec.Version] = pv
		}
	}

	verConfig := versions.VersionSelection{Semver: semverConfig}

	selectedVersion, err := versions.HighestConstrainedVersion(versionStrs, verConfig)
	if err != nil {
		return datapkgingv1alpha1.Package{}, err
	}

	if pkg, found := versionToPkg[selectedVersion]; found {
		return pkg, nil
	}

	return datapkgingv1alpha1.Package{}, fmt.Errorf("Could not find package with name '%s' and version '%s'",
		pi.model.Spec.PackageRef.RefName, selectedVersion)
}

func (pi *PackageInstallCR) reconcileDelete(modelStatus *reconciler.Status) (reconcile.Result, error) {
	existingApp, err := pi.kcclient.KappctrlV1alpha1().Apps(pi.model.Namespace).Get(
		context.Background(), pi.model.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, pi.unblockDeletion()
		}
		return reconcile.Result{Requeue: true}, err
	}

	if existingApp.DeletionTimestamp == nil {
		err := pi.kcclient.KappctrlV1alpha1().Apps(pi.model.Namespace).Delete(
			context.Background(), pi.model.Name, metav1.DeleteOptions{})
		return reconcile.Result{}, err
	}

	appStatus := reconciler.Status{S: existingApp.Status.GenericStatus}
	pi.log.Info(fmt.Sprintf("Reconciling deletion of PackageInstall '%s/%s'", pi.model.Namespace, pi.model.Name))
	switch {
	case appStatus.IsDeleting():
		modelStatus.SetDeleting(pi.model.ObjectMeta)
	case appStatus.IsDeleteFailed():
		modelStatus.SetUsefulErrorMessage(existingApp.Status.UsefulErrorMessage)
		modelStatus.SetDeleteCompleted(fmt.Errorf("Error (see .status.usefulErrorMessage for details)"))
	}

	return reconcile.Result{}, nil // Nothing to do
}

func (pi *PackageInstallCR) updateStatus() error {
	if !equality.Semantic.DeepEqual(pi.unmodifiedModel.Status, pi.model.Status) {
		_, err := pi.kcclient.PackagingV1alpha1().PackageInstalls(pi.model.Namespace).UpdateStatus(context.Background(), pi.model, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("Updating installed package status: %s", err)
		}
	}
	return nil
}

func (pi *PackageInstallCR) blockDeletion() error {
	// Avoid doing unnecessary processing
	if containsString(pi.unmodifiedModel.Finalizers, deleteFinalizerName) {
		return nil
	}

	pi.log.Info("Blocking deletion")

	return pi.update(func(ipkg *pkgingv1alpha1.PackageInstall) {
		if !containsString(ipkg.ObjectMeta.Finalizers, deleteFinalizerName) {
			ipkg.ObjectMeta.Finalizers = append(ipkg.ObjectMeta.Finalizers, deleteFinalizerName)
		}
	})
}

func (pi *PackageInstallCR) unblockDeletion() error {
	pi.log.Info("Unblocking deletion")
	return pi.update(func(ipkg *pkgingv1alpha1.PackageInstall) {
		ipkg.ObjectMeta.Finalizers = removeString(ipkg.ObjectMeta.Finalizers, deleteFinalizerName)
	})
}

func (pi *PackageInstallCR) update(updateFunc func(*pkgingv1alpha1.PackageInstall)) error {
	pi.log.Info("Updating installed package")

	modelForUpdate := pi.model.DeepCopy()
	updateFunc(modelForUpdate)

	var err error

	pi.model, err = pi.kcclient.PackagingV1alpha1().PackageInstalls(modelForUpdate.Namespace).Update(
		context.Background(), modelForUpdate, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Updating installed package: %s", err)
	}

	pi.unmodifiedModel = pi.model.DeepCopy()

	return nil
}
