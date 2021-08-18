// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	pkgclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reconciler"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

	pv, err := pi.referencedPkgVersion()
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	pi.model.Status.Version = pv.Spec.Version

	existingApp, err := pi.kcclient.KappctrlV1alpha1().Apps(pi.model.Namespace).Get(context.Background(), pi.model.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			err := pi.reconcileFetchPlaceholderSecrets(&pv)
			if err != nil {
				return reconcile.Result{}, err
			}
			return pi.createAppFromPackage(pv)
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
	err := pi.reconcileFetchPlaceholderSecrets(&pv)
	if err != nil {
		return reconcile.Result{}, err
	}

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
	pi.log.Info("Reconciling deletion")

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

func (pi *PackageInstallCR) reconcileFetchPlaceholderSecrets(pv *datapkgingv1alpha1.Package) error {
	for i, fetch := range pv.Spec.Template.Spec.Fetch {
		if fetch.ImgpkgBundle != nil && fetch.ImgpkgBundle.SecretRef == nil {
			secretName, err := pi.createSecretForSecretgenController(i)
			if err != nil {
				return err
			}
			pv.Spec.Template.Spec.Fetch[i].ImgpkgBundle.SecretRef = &kcv1alpha1.AppFetchLocalRef{secretName}
		}

		if fetch.Image != nil && fetch.Image.SecretRef == nil {
			secretName, err := pi.createSecretForSecretgenController(i)
			if err != nil {
				return err
			}
			pv.Spec.Template.Spec.Fetch[i].Image.SecretRef = &kcv1alpha1.AppFetchLocalRef{secretName}
		}
	}
	return nil
}

func (pi PackageInstallCR) createSecretForSecretgenController(iteration int) (string, error) {
	secretName := fmt.Sprintf("%s-%s", pi.model.Name, "fetch"+strconv.Itoa(iteration))
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: pi.model.Namespace,
			Annotations: map[string]string{
				"secretgen.carvel.dev/image-pull-secret": "",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "packaging.carvel.dev/v1alpha1",
					Kind:       "PackageInstall",
					Name:       pi.model.Name,
					UID:        pi.model.UID,
				},
			},
		},
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: []byte(`{"auths":{}}`),
		},
		Type: corev1.SecretTypeDockerConfigJson,
	}

	_, err := pi.coreClient.CoreV1().Secrets(pi.model.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			_, err = pi.coreClient.CoreV1().Secrets(pi.model.Namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return secretName, nil
}
