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
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/scheme"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reconciler"
	"github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions"
	verv1alpha1 "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// DowngradableAnnKey specifies annotation that user can place on
	// PackageInstall to indicate that lower version of the package
	// can be selected vs whats currently installed.
	DowngradableAnnKey = "packaging.carvel.dev/downgradable"
)

// nolint: revive
type PackageInstallCR struct {
	model           *pkgingv1alpha1.PackageInstall
	unmodifiedModel *pkgingv1alpha1.PackageInstall

	log        logr.Logger
	kcclient   kcclient.Interface
	pkgclient  pkgclient.Interface
	coreClient kubernetes.Interface
}

func NewPackageInstallCR(model *pkgingv1alpha1.PackageInstall, log logr.Logger,
	kcclient kcclient.Interface, pkgclient pkgclient.Interface, coreClient kubernetes.Interface) *PackageInstallCR {

	return &PackageInstallCR{model: model, unmodifiedModel: model.DeepCopy(), log: log,
		kcclient: kcclient, pkgclient: pkgclient, coreClient: coreClient}
}

func (pi *PackageInstallCR) Reconcile() (reconcile.Result, error) {
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
	pi.log.Info("Reconciling")

	err := pi.blockDeletion()
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	modelStatus.SetReconciling(pi.model.ObjectMeta)

	var fieldErrors field.ErrorList
	for i, value := range pi.model.Spec.Values {
		if value.SecretRef == nil {
			fieldErrors = append(fieldErrors, field.Required(field.NewPath("spec", "values").Index(i).Child("secretRef"), ""))
		}
	}
	if len(fieldErrors) > 0 {
		return reconcile.Result{}, fmt.Errorf("Invalid fields: %w", fieldErrors.ToAggregate())
	}
	pkg, err := pi.referencedPkgVersion()
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	// Set new desired version before checking if it's not applicable
	pi.model.Status.Version = pkg.Spec.Version

	_, canDowngrade := pi.model.Annotations[DowngradableAnnKey]
	if !canDowngrade && pi.model.Status.LastAttemptedVersion != "" {
		matchedVers := versions.NewRelaxedSemversNoErr([]string{pkg.Spec.Version})

		matchedVers, err = matchedVers.FilterConstraints(">=" + pi.model.Status.LastAttemptedVersion)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("Filtering by last attempted version '%s': %s",
				pi.model.Status.LastAttemptedVersion, err)
		}

		if matchedVers.Len() == 0 {
			errMsg := fmt.Sprintf(
				"Stopped installing matched version '%s' since last attempted version '%s' is higher."+
					"\nhint: Add annotation packaging.carvel.dev/downgradable: \"\" to PackageInstall to proceed with downgrade",
				pkg.Spec.Version, pi.model.Status.LastAttemptedVersion)
			modelStatus.SetUsefulErrorMessage(errMsg)
			modelStatus.SetReconcileCompleted(fmt.Errorf("Error (see .status.usefulErrorMessage for details)"))
			// Nothing to do until available packages change or PackageInstall changes
			return reconcile.Result{}, nil
		}
	}

	pi.model.Status.LastAttemptedVersion = pkg.Spec.Version

	existingApp, err := pi.kcclient.KappctrlV1alpha1().Apps(pi.model.Namespace).Get(
		context.Background(), pi.model.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			pkgWithPlaceholderSecrets, err := pi.reconcileFetchPlaceholderSecrets(pkg)
			if err != nil {
				return reconcile.Result{}, err
			}
			return pi.createAppFromPackage(pkgWithPlaceholderSecrets)
		}
		return reconcile.Result{Requeue: true}, err
	}

	appStatus := reconciler.Status{S: existingApp.Status.GenericStatus}
	switch {
	case appStatus.IsReconciling():
		modelStatus.SetReconciling(pi.model.ObjectMeta)
	case appStatus.IsReconcileSucceeded():
		modelStatus.SetReconcileCompleted(nil)
	case appStatus.IsReconcileFailed():
		modelStatus.SetUsefulErrorMessage(existingApp.Status.UsefulErrorMessage)
		modelStatus.SetReconcileCompleted(fmt.Errorf("Error (see .status.usefulErrorMessage for details)"))
	}

	return pi.reconcileAppWithPackage(existingApp, pkg)
}

func (pi *PackageInstallCR) createAppFromPackage(pkg datapkgingv1alpha1.Package) (reconcile.Result, error) {
	desiredApp, err := NewApp(&v1alpha1.App{}, pi.model, pkg)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	_, err = pi.kcclient.KappctrlV1alpha1().Apps(desiredApp.Namespace).Create(context.Background(), desiredApp, metav1.CreateOptions{})
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	return reconcile.Result{}, nil
}

func (pi *PackageInstallCR) reconcileAppWithPackage(existingApp *kcv1alpha1.App, pkg datapkgingv1alpha1.Package) (reconcile.Result, error) {
	pkgWithPlaceholderSecrets, err := pi.reconcileFetchPlaceholderSecrets(pkg)
	if err != nil {
		return reconcile.Result{}, err
	}

	desiredApp, err := NewApp(existingApp, pi.model, pkgWithPlaceholderSecrets)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	if !equality.Semantic.DeepEqual(desiredApp, existingApp) {
		_, err = pi.kcclient.KappctrlV1alpha1().Apps(desiredApp.Namespace).Update(
			context.Background(), desiredApp, metav1.UpdateOptions{})
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

	pkgList, err := pi.pkgclient.DataV1alpha1().Packages(pi.model.Namespace).List(
		context.Background(), metav1.ListOptions{})
	if err != nil {
		return datapkgingv1alpha1.Package{}, err
	}

	var versionStrs []string
	versionToPkg := map[string]datapkgingv1alpha1.Package{}

	for _, pkg := range pkgList.Items {
		if pkg.Spec.RefName == pi.model.Spec.PackageRef.RefName {
			versionStrs = append(versionStrs, pkg.Spec.Version)
			versionToPkg[pkg.Spec.Version] = pkg
		}
	}

	// If constraint is a single specified version, then we
	// do not want to force user to manually set prereleases={}
	if len(semverConfig.Constraints) > 0 && semverConfig.Prereleases == nil {
		// Will error if it's not a single version
		singleVer, err := versions.NewSemver(semverConfig.Constraints)
		if err == nil && len(singleVer.Pre) > 0 {
			semverConfig.Prereleases = &verv1alpha1.VersionSelectionSemverPrereleases{}
		}
	}

	verConfig := verv1alpha1.VersionSelection{Semver: semverConfig}

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
	pi.log.Info("Reconciling deletion")

	existingApp, err := pi.kcclient.KappctrlV1alpha1().Apps(pi.model.Namespace).Get(
		context.Background(), pi.model.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, pi.unblockDeletion()
		}
		return reconcile.Result{Requeue: true}, err
	}

	unchangeExistingApp := existingApp.DeepCopy()

	// Ensure that several fields that may affect how App is deleted
	// are set to same values as they are on PackageInstall
	if existingApp.Spec.ServiceAccountName != pi.model.Spec.ServiceAccountName {
		existingApp.Spec.ServiceAccountName = pi.model.Spec.ServiceAccountName
	}
	if existingApp.Spec.Cluster != pi.model.Spec.Cluster {
		existingApp.Spec.Cluster = pi.model.Spec.Cluster
	}
	if existingApp.Spec.NoopDelete != pi.model.Spec.NoopDelete {
		existingApp.Spec.NoopDelete = pi.model.Spec.NoopDelete
	}
	if existingApp.Spec.Paused != pi.model.Spec.Paused {
		existingApp.Spec.Paused = pi.model.Spec.Paused
	}
	if existingApp.Spec.Canceled != pi.model.Spec.Canceled {
		existingApp.Spec.Canceled = pi.model.Spec.Canceled
	}

	if !equality.Semantic.DeepEqual(existingApp, unchangeExistingApp) {
		existingApp, err = pi.kcclient.KappctrlV1alpha1().Apps(existingApp.Namespace).Update(
			context.Background(), existingApp, metav1.UpdateOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return reconcile.Result{}, pi.unblockDeletion()
			}
			return reconcile.Result{Requeue: true}, err
		}
	}

	if existingApp.DeletionTimestamp == nil {
		err := pi.kcclient.KappctrlV1alpha1().Apps(existingApp.Namespace).Delete(
			context.Background(), existingApp.Name, metav1.DeleteOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return reconcile.Result{}, pi.unblockDeletion()
			}
		}
		return reconcile.Result{}, err
	}

	appStatus := reconciler.Status{S: existingApp.Status.GenericStatus}
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
		_, err := pi.kcclient.PackagingV1alpha1().PackageInstalls(pi.model.Namespace).UpdateStatus(
			context.Background(), pi.model, metav1.UpdateOptions{})
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

	var lastErr error
	for i := 0; i < 5; i++ {
		updateFunc(modelForUpdate)

		updatedModel, err := pi.kcclient.PackagingV1alpha1().PackageInstalls(modelForUpdate.Namespace).Update(
			context.Background(), modelForUpdate, metav1.UpdateOptions{})
		if err == nil {
			pi.model = updatedModel
			pi.unmodifiedModel = updatedModel.DeepCopy()
			return nil
		}

		lastErr = err

		// if we errored, refresh the model we have
		modelForUpdate, err = pi.kcclient.PackagingV1alpha1().PackageInstalls(modelForUpdate.Namespace).Get(
			context.Background(), modelForUpdate.Name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("Getting package install model: %s", err)
		}
	}

	return fmt.Errorf("Updating package install: %s", lastErr)
}

func (pi *PackageInstallCR) reconcileFetchPlaceholderSecrets(pkg datapkgingv1alpha1.Package) (datapkgingv1alpha1.Package, error) {
	pkg = *pkg.DeepCopy()
	for i, fetch := range pkg.Spec.Template.Spec.Fetch {
		if fetch.ImgpkgBundle != nil && fetch.ImgpkgBundle.SecretRef == nil {
			secretName, err := pi.createSecretForSecretgenController(i)
			if err != nil {
				return datapkgingv1alpha1.Package{}, err
			}
			pkg.Spec.Template.Spec.Fetch[i].ImgpkgBundle.SecretRef = &kcv1alpha1.AppFetchLocalRef{secretName}
		}

		if fetch.Image != nil && fetch.Image.SecretRef == nil {
			secretName, err := pi.createSecretForSecretgenController(i)
			if err != nil {
				return datapkgingv1alpha1.Package{}, err
			}
			pkg.Spec.Template.Spec.Fetch[i].Image.SecretRef = &kcv1alpha1.AppFetchLocalRef{secretName}
		}
	}
	return pkg, nil
}

func (pi PackageInstallCR) createSecretForSecretgenController(iteration int) (string, error) {
	secretName := fmt.Sprintf("%s-fetch-%d", pi.model.Name, iteration)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: pi.model.Namespace,
			Annotations: map[string]string{
				"secretgen.carvel.dev/image-pull-secret": "",
			},
		},
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: []byte(`{"auths":{}}`),
		},
		Type: corev1.SecretTypeDockerConfigJson,
	}

	controllerutil.SetOwnerReference(pi.model, secret, scheme.Scheme)

	_, err := pi.coreClient.CoreV1().Secrets(pi.model.Namespace).Create(
		context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return "", err
		}
	}
	return secretName, nil
}
