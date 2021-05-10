// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installedpkg

import (
	"context"
	"fmt"

	instPkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	pkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/packages/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reconciler"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type InstalledPackageCR struct {
	model           *instPkgv1alpha1.InstalledPackage
	unmodifiedModel *instPkgv1alpha1.InstalledPackage

	log       logr.Logger
	kcclient  kcclient.Interface
	pkgclient pkgclient.Interface
}

func NewInstalledPkgCR(model *instPkgv1alpha1.InstalledPackage, log logr.Logger,
	kcclient kcclient.Interface, pkgclient pkgclient.Interface) *InstalledPackageCR {

	return &InstalledPackageCR{model: model, unmodifiedModel: model.DeepCopy(), log: log, kcclient: kcclient, pkgclient: pkgclient}
}

func (ip *InstalledPackageCR) Reconcile() (reconcile.Result, error) {
	ip.log.Info(fmt.Sprintf("Reconciling InstalledPackage '%s/%s'", ip.model.Namespace, ip.model.Name))

	// TODO deleting conditions
	if ip.model.DeletionTimestamp != nil {
		return reconcile.Result{}, nil // Nothing to do
	}

	status := &reconciler.Status{
		ip.model.Status.GenericStatus,
		func(st kcv1alpha1.GenericStatus) { ip.model.Status.GenericStatus = st },
	}

	result, err := ip.reconcile(status)
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

func (ip *InstalledPackageCR) reconcile(modelStatus *reconciler.Status) (reconcile.Result, error) {
	modelStatus.SetReconciling(ip.model.ObjectMeta)

	pv, err := ip.referencedPkgVersion()
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	ip.model.Status.Version = pv.Spec.Version

	existingApp, err := ip.kcclient.KappctrlV1alpha1().Apps(ip.model.Namespace).Get(context.Background(), ip.model.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return ip.createAppFromPackage(pv)
		}
		return reconcile.Result{Requeue: true}, err
	}

	appStatus := reconciler.Status{S: existingApp.Status.GenericStatus}
	ip.log.Info(fmt.Sprintf("Reconciling InstalledPackage '%s/%s'", ip.model.Namespace, ip.model.Name))
	switch {
	case appStatus.IsReconciling():
		modelStatus.SetReconciling(ip.model.ObjectMeta)
	case appStatus.IsReconcileSucceeded():
		modelStatus.SetReconcileCompleted(nil)
	case appStatus.IsReconcileFailed():
		modelStatus.SetUsefulErrorMessage(existingApp.Status.UsefulErrorMessage)
		modelStatus.SetReconcileCompleted(fmt.Errorf("Error (see .status.usefulErrorMessage for details)"))
	}

	return ip.reconcileAppWithPackage(existingApp, pv)
}

func (ip *InstalledPackageCR) createAppFromPackage(pv pkgv1alpha1.PackageVersion) (reconcile.Result, error) {
	desiredApp, err := NewApp(&v1alpha1.App{}, ip.model, pv)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	_, err = ip.kcclient.KappctrlV1alpha1().Apps(desiredApp.Namespace).Create(context.Background(), desiredApp, metav1.CreateOptions{})
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	return reconcile.Result{}, nil
}

func (ip *InstalledPackageCR) reconcileAppWithPackage(existingApp *kcv1alpha1.App, pv pkgv1alpha1.PackageVersion) (reconcile.Result, error) {
	desiredApp, err := NewApp(existingApp, ip.model, pv)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	if !equality.Semantic.DeepEqual(desiredApp, existingApp) {
		_, err = ip.kcclient.KappctrlV1alpha1().Apps(desiredApp.Namespace).Update(context.Background(), desiredApp, metav1.UpdateOptions{})
		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}
	}

	return reconcile.Result{}, nil
}

func (ip *InstalledPackageCR) referencedPkgVersion() (pkgv1alpha1.PackageVersion, error) {
	if ip.model.Spec.PackageVersionRef == nil {
		return pkgv1alpha1.PackageVersion{}, fmt.Errorf("Expected non nil PackageVersionRef")
	}

	semverConfig := ip.model.Spec.PackageVersionRef.VersionSelection

	pvList, err := ip.pkgclient.PackageV1alpha1().PackageVersions().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return pkgv1alpha1.PackageVersion{}, err
	}

	var versionStrs []string
	versionToPkg := map[string]pkgv1alpha1.PackageVersion{}

	for _, pv := range pvList.Items {
		if pv.Spec.PackageName == ip.model.Spec.PackageVersionRef.PackageName {
			versionStrs = append(versionStrs, pv.Spec.Version)
			versionToPkg[pv.Spec.Version] = pv
		}
	}

	verConfig := versions.VersionSelection{Semver: semverConfig}

	selectedVersion, err := versions.HighestConstrainedVersion(versionStrs, verConfig)
	if err != nil {
		return pkgv1alpha1.PackageVersion{}, err
	}

	if pkg, found := versionToPkg[selectedVersion]; found {
		return pkg, nil
	}

	return pkgv1alpha1.PackageVersion{}, fmt.Errorf("Could not find package with name '%s' and version '%s'",
		ip.model.Spec.PackageVersionRef.PackageName, selectedVersion)
}

func (ip *InstalledPackageCR) updateStatus() error {
	if !equality.Semantic.DeepEqual(ip.unmodifiedModel.Status, ip.model.Status) {
		_, err := ip.kcclient.InstallV1alpha1().InstalledPackages(ip.model.Namespace).UpdateStatus(context.Background(), ip.model, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("Updating installed package status: %s", err)
		}
	}
	return nil
}
