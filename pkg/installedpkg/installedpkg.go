// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installedpkg

import (
	"fmt"
	instPkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	pkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/package/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reconciler"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type InstalledPackageCR struct {
	model           *instPkgv1alpha1.InstalledPackage
	unmodifiedModel *instPkgv1alpha1.InstalledPackage

	log    logr.Logger
	client kcclient.Interface
}

func init() {
	instPkgv1alpha1.AddToScheme(scheme.Scheme)
	pkgv1alpha1.AddToScheme(scheme.Scheme)
}

func NewInstalledPkgCR(model *instPkgv1alpha1.InstalledPackage, log logr.Logger,
	client kcclient.Interface) *InstalledPackageCR {

	return &InstalledPackageCR{model: model, unmodifiedModel: model.DeepCopy(), log: log, client: client}
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

	pkg, err := ip.referencedPkg()
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	ip.model.Status.Version = pkg.Spec.Version

	existingApp, err := ip.client.KappctrlV1alpha1().Apps(ip.model.Namespace).Get(ip.model.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return ip.createAppFromPackage(pkg)
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
		modelStatus.SetReconcileCompleted(fmt.Errorf("App failed reconciling"))
	}

	return ip.reconcileAppWithPackage(existingApp, pkg)
}

func (ip *InstalledPackageCR) createAppFromPackage(pkg pkgv1alpha1.Package) (reconcile.Result, error) {
	desiredApp, err := NewApp(&v1alpha1.App{}, ip.model, pkg)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	_, err = ip.client.KappctrlV1alpha1().Apps(desiredApp.Namespace).Create(desiredApp)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	return reconcile.Result{}, nil
}

func (ip *InstalledPackageCR) reconcileAppWithPackage(existingApp *kcv1alpha1.App, pkg pkgv1alpha1.Package) (reconcile.Result, error) {
	desiredApp, err := NewApp(existingApp, ip.model, pkg)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	if !equality.Semantic.DeepEqual(desiredApp, existingApp) {
		_, err = ip.client.KappctrlV1alpha1().Apps(desiredApp.Namespace).Update(desiredApp)
		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}
	}

	return reconcile.Result{}, nil
}

func (ip *InstalledPackageCR) referencedPkg() (pkgv1alpha1.Package, error) {
	var semverConfig *versions.VersionSelectionSemver

	switch {
	case ip.model.Spec.PkgRef.Version != "" && ip.model.Spec.PkgRef.VersionSelection != nil:
		return pkgv1alpha1.Package{}, fmt.Errorf("Cannot use 'version' with 'versionSelection'")
	case ip.model.Spec.PkgRef.Version != "":
		semverConfig = &versions.VersionSelectionSemver{Constraints: ip.model.Spec.PkgRef.Version}
	case ip.model.Spec.PkgRef.VersionSelection != nil:
		semverConfig = ip.model.Spec.PkgRef.VersionSelection
	}

	pkgList, err := ip.client.PackageV1alpha1().Packages().List(metav1.ListOptions{})
	if err != nil {
		return pkgv1alpha1.Package{}, err
	}

	var versionStrs []string
	versionToPkg := map[string]pkgv1alpha1.Package{}

	for _, pkg := range pkgList.Items {
		versionStrs = append(versionStrs, pkg.Spec.Version)
		versionToPkg[pkg.Spec.Version] = pkg
	}

	verConfig := versions.VersionSelection{Semver: semverConfig}

	selectedVersion, err := versions.HighestConstrainedVersion(versionStrs, verConfig)
	if err != nil {
		return pkgv1alpha1.Package{}, err
	}

	if pkg, found := versionToPkg[selectedVersion]; found {
		return pkg, nil
	}

	return pkgv1alpha1.Package{}, fmt.Errorf("Could not find package with name '%s' and version '%s'",
		ip.model.Spec.PkgRef.PublicName, selectedVersion)
}

func (ip *InstalledPackageCR) updateStatus() error {
	if !equality.Semantic.DeepEqual(ip.unmodifiedModel.Status, ip.model.Status) {
		_, err := ip.client.InstallV1alpha1().InstalledPackages(ip.model.Namespace).UpdateStatus(ip.model)
		if err != nil {
			return fmt.Errorf("Updating installed package status: %s", err)
		}
	}
	return nil
}
