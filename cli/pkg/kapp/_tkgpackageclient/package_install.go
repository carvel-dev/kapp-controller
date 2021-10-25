// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"

	kappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

const (
	msgRunPackageInstalledDelete = "\n\nPlease consider using 'tanzu package installed delete' to delete the already created associated resources\n"
	msgRunPackageInstalledUpdate = "\n\nPlease consider using 'tanzu package installed update' to update the installed package with correct settings\n"
)

// InstallPackage installs the PackageInstall and its associated resources in the cluster
func (p *pkgClient) InstallPackage(o *tkgpackagedatamodel.PackageOptions, progress *tkgpackagedatamodel.PackageProgress, operationType tkgpackagedatamodel.OperationType) { //nolint:gocyclo
	var (
		pkgInstall            *kappipkg.PackageInstall
		err                   error
		secretCreated         bool
		serviceAccountCreated bool
	)

	defer func() {
		if err != nil {
			progress.Err <- err
		}
		if operationType == tkgpackagedatamodel.OperationTypeInstall {
			close(progress.ProgressMsg)
			close(progress.Done)
		}
	}()

	if pkgInstall, err = p.kappClient.GetPackageInstall(o.PkgInstallName, o.Namespace); err != nil {
		if !k8serror.IsNotFound(err) {
			return
		}
		err = nil
	}

	if pkgInstall != nil && pkgInstall.Name == o.PkgInstallName {
		err = errors.New(fmt.Sprintf("package install '%s' already exists in namespace '%s'", o.PkgInstallName, o.Namespace))
		return
	}

	progress.ProgressMsg <- fmt.Sprintf("Getting package metadata for '%s'", o.PackageName)
	if _, _, err = p.GetPackage(o); err != nil {
		return
	}

	if o.CreateNamespace {
		progress.ProgressMsg <- fmt.Sprintf("Creating namespace '%s'", o.Namespace)
		if err = p.createNamespace(o.Namespace); err != nil {
			return
		}
	} else if err = p.kappClient.GetClient().Get(context.Background(), crtclient.ObjectKey{Name: o.Namespace}, &corev1.Namespace{}); err != nil {
		return
	}

	if o.ServiceAccountName == "" {
		o.ServiceAccountName = fmt.Sprintf(tkgpackagedatamodel.ServiceAccountName, o.PkgInstallName, o.Namespace)
		progress.ProgressMsg <- fmt.Sprintf("Creating service account '%s'", o.ServiceAccountName)
		if serviceAccountCreated, err = p.createServiceAccount(o); err != nil {
			return
		}

		o.ClusterRoleName = fmt.Sprintf(tkgpackagedatamodel.ClusterRoleName, o.PkgInstallName, o.Namespace)
		progress.ProgressMsg <- fmt.Sprintf("Creating cluster admin role '%s'", o.ClusterRoleName)
		if err = p.createClusterAdminRole(o); err != nil {
			log.Warning(msgRunPackageInstalledDelete)
			return
		}

		o.ClusterRoleBindingName = fmt.Sprintf(tkgpackagedatamodel.ClusterRoleBindingName, o.PkgInstallName, o.Namespace)
		progress.ProgressMsg <- fmt.Sprintf("Creating cluster role binding '%s'", o.ClusterRoleBindingName)
		if err = p.createClusterRoleBinding(o); err != nil {
			log.Warning(msgRunPackageInstalledDelete)
			return
		}
	} else {
		objKey := crtclient.ObjectKey{Name: o.ServiceAccountName, Namespace: o.Namespace}
		svcAccount := &corev1.ServiceAccount{}
		if err = p.kappClient.GetClient().Get(context.Background(), objKey, svcAccount); err != nil {
			err = errors.Wrap(err, fmt.Sprintf("failed to find service account '%s' in namespace '%s'", o.ServiceAccountName, o.Namespace))
			return
		}
		if _, ok := svcAccount.GetAnnotations()[tkgpackagedatamodel.TanzuPkgPluginAnnotation]; ok {
			err = errors.New(fmt.Sprintf("provided service account '%s' is already used by another package in namespace '%s'", o.ServiceAccountName, o.Namespace))
			return
		}
	}

	if o.ValuesFile != "" {
		o.SecretName = fmt.Sprintf(tkgpackagedatamodel.SecretName, o.PkgInstallName, o.Namespace)
		progress.ProgressMsg <- fmt.Sprintf("Creating secret '%s'", o.SecretName)
		if secretCreated, err = p.createDataValuesSecret(o); err != nil {
			log.Warning(msgRunPackageInstalledDelete)
			return
		}
	}

	progress.ProgressMsg <- "Creating package resource"
	if err = p.createPackageInstall(o, serviceAccountCreated, secretCreated); err != nil {
		log.Warning(msgRunPackageInstalledDelete)
		return
	}

	if o.Wait {
		if err = p.waitForResourceInstallation(o.PkgInstallName, o.Namespace, o.PollInterval, o.PollTimeout, progress.ProgressMsg, tkgpackagedatamodel.ResourceTypePackageInstall); err != nil {
			log.Warning(msgRunPackageInstalledUpdate)
			return
		}
	}
}

// createClusterAdminRole creates a ClusterRole resource
func (p *pkgClient) createClusterAdminRole(o *tkgpackagedatamodel.PackageOptions) error {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:        o.ClusterRoleName,
			Annotations: map[string]string{tkgpackagedatamodel.TanzuPkgPluginAnnotation: fmt.Sprintf(tkgpackagedatamodel.TanzuPkgPluginResource, o.PkgInstallName, o.Namespace)},
		},
		Rules: []rbacv1.PolicyRule{
			{APIGroups: []string{"*"}, Verbs: []string{"*"}, Resources: []string{"*"}},
		},
	}

	if err := p.kappClient.GetClient().Create(context.Background(), clusterRole); err != nil {
		return errors.Wrap(err, "failed to create ClusterRole resource")
	}

	return nil
}

// createClusterRoleBinding creates a ClusterRoleBinding resource
func (p *pkgClient) createClusterRoleBinding(o *tkgpackagedatamodel.PackageOptions) error {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        o.ClusterRoleBindingName,
			Annotations: map[string]string{tkgpackagedatamodel.TanzuPkgPluginAnnotation: fmt.Sprintf(tkgpackagedatamodel.TanzuPkgPluginResource, o.PkgInstallName, o.Namespace)},
		},
		Subjects: []rbacv1.Subject{{Kind: tkgpackagedatamodel.KindServiceAccount, Name: o.ServiceAccountName, Namespace: o.Namespace}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     tkgpackagedatamodel.KindClusterRole,
			Name:     o.ClusterRoleName,
		},
	}

	if err := p.kappClient.GetClient().Create(context.Background(), clusterRoleBinding); err != nil {
		return errors.Wrap(err, "failed to create ClusterRoleBinding resource")
	}

	return nil
}

// createDataValuesSecret create a secret object containing the user-provided configuration.
func (p *pkgClient) createDataValuesSecret(o *tkgpackagedatamodel.PackageOptions) (bool, error) {
	var err error

	dataValues := make(map[string][]byte)

	if dataValues[filepath.Base(o.ValuesFile)], err = ioutil.ReadFile(o.ValuesFile); err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("failed to read from data values file '%s'", o.ValuesFile))
	}
	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        o.SecretName,
			Namespace:   o.Namespace,
			Annotations: map[string]string{tkgpackagedatamodel.TanzuPkgPluginAnnotation: fmt.Sprintf(tkgpackagedatamodel.TanzuPkgPluginResource, o.PkgInstallName, o.Namespace)},
		},
		Data: dataValues,
	}

	if err := p.kappClient.GetClient().Create(context.Background(), secret); err != nil {
		return false, errors.Wrap(err, "failed to create Secret resource")
	}

	return true, nil
}

// createNamespace creates a namespace resource if it doesn't already exist
func (p *pkgClient) createNamespace(namespace string) error {
	err := p.kappClient.GetClient().Get(
		context.Background(),
		crtclient.ObjectKey{Name: namespace},
		&corev1.Namespace{})
	if err != nil {
		if !k8serror.IsNotFound(err) {
			return err
		}
		ns := &corev1.Namespace{
			TypeMeta:   metav1.TypeMeta{Kind: tkgpackagedatamodel.KindNamespace},
			ObjectMeta: metav1.ObjectMeta{Name: namespace},
		}
		if err := p.kappClient.GetClient().Create(context.Background(), ns); err != nil {
			return errors.Wrap(err, "failed to create namespace")
		}
	}

	return nil
}

// createPackageInstall creates the PackageInstall CR
func (p *pkgClient) createPackageInstall(o *tkgpackagedatamodel.PackageOptions, serviceAccountCreated, secretCreated bool) error {
	// construct the PackageInstall CR
	packageInstall := &kappipkg.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{Name: o.PkgInstallName, Namespace: o.Namespace},
		Spec: kappipkg.PackageInstallSpec{
			ServiceAccountName: o.ServiceAccountName,
			PackageRef: &kappipkg.PackageRef{
				RefName: o.PackageName,
				VersionSelection: &versions.VersionSelectionSemver{
					Constraints: o.Version,
					Prereleases: &versions.VersionSelectionSemverPrereleases{},
				},
			},
		},
	}

	// if configuration data file was provided, reference the secret name in the PackageInstall
	if secretCreated {
		packageInstall.Spec.Values = []kappipkg.PackageInstallValues{
			{
				SecretRef: &kappipkg.PackageInstallValuesSecretRef{
					Name: fmt.Sprintf(tkgpackagedatamodel.SecretName, o.PkgInstallName, o.Namespace),
				},
			},
		}
	}

	if err := p.kappClient.CreatePackageInstall(packageInstall, serviceAccountCreated, secretCreated); err != nil {
		return errors.Wrap(err, "failed to create PackageInstall resource")
	}

	return nil
}

// createServiceAccount creates a ServiceAccount resource
func (p *pkgClient) createServiceAccount(o *tkgpackagedatamodel.PackageOptions) (bool, error) {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        o.ServiceAccountName,
			Namespace:   o.Namespace,
			Annotations: map[string]string{tkgpackagedatamodel.TanzuPkgPluginAnnotation: fmt.Sprintf(tkgpackagedatamodel.TanzuPkgPluginResource, o.PkgInstallName, o.Namespace)}},
	}

	if err := p.kappClient.GetClient().Create(context.Background(), serviceAccount); err != nil {
		return false, errors.Wrap(err, "failed to create ServiceAccount resource")
	}

	return true, nil
}

// waitForResourceInstallation waits until the package get installed successfully or a failure happen
func (p *pkgClient) waitForResourceInstallation(name, namespace string, pollInterval, pollTimeout time.Duration, progress chan string, rscType tkgpackagedatamodel.ResourceType) error { //nolint:gocyclo
	var (
		status             kappctrl.GenericStatus
		reconcileSucceeded bool
	)
	progress <- fmt.Sprintf("Waiting for '%s' reconciliation for '%s'", rscType.String(), name)
	if err := wait.Poll(pollInterval, pollTimeout, func() (done bool, err error) {
		switch rscType {
		case tkgpackagedatamodel.ResourceTypePackageRepository:
			resource, err := p.kappClient.GetPackageRepository(name, namespace)
			if err != nil {
				return false, err
			}
			if resource.Generation != resource.Status.ObservedGeneration {
				// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
				return false, nil
			}
			status = resource.Status.GenericStatus
		case tkgpackagedatamodel.ResourceTypePackageInstall:
			resource, err := p.kappClient.GetPackageInstall(name, namespace)
			if err != nil {
				return false, err
			}
			if resource.Generation != resource.Status.ObservedGeneration {
				// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
				return false, nil
			}
			status = resource.Status.GenericStatus
		}

		for _, cond := range status.Conditions {
			if progress != nil {
				progress <- fmt.Sprintf("'%s' resource install status: %s", rscType.String(), cond.Type)
			}
			switch {
			case cond.Type == kappctrl.ReconcileSucceeded && cond.Status == corev1.ConditionTrue:
				if progress != nil {
					progress <- fmt.Sprintf("'%s' resource successfully reconciled", rscType.String())
				}
				reconcileSucceeded = true
				return true, nil
			case cond.Type == kappctrl.ReconcileFailed && cond.Status == corev1.ConditionTrue:
				return false, fmt.Errorf("resource reconciliation failed: %s. %s", status.UsefulErrorMessage, status.FriendlyDescription)
			}
		}
		return false, nil
	}); err != nil {
		return err
	}

	if !reconcileSucceeded {
		return errors.New(fmt.Sprintf("'%s' resource reconciliation failed", rscType.String()))
	}

	return nil
}
