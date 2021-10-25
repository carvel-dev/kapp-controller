// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

func (p *pkgClient) UpdatePackage(o *tkgpackagedatamodel.PackageOptions, progress *tkgpackagedatamodel.PackageProgress) {
	var (
		pkgInstall    *kappipkg.PackageInstall
		err           error
		secretCreated bool
	)

	defer func() {
		progressCleanup(err, progress)
	}()

	progress.ProgressMsg <- fmt.Sprintf("Getting package install for '%s'", o.PkgInstallName)
	if pkgInstall, err = p.kappClient.GetPackageInstall(o.PkgInstallName, o.Namespace); err != nil {
		if !k8serror.IsNotFound(err) {
			return
		}
		err = nil
	}

	if pkgInstall == nil {
		if !o.Install {
			err = &tkgpackagedatamodel.PackagePluginNonCriticalError{Reason: tkgpackagedatamodel.ErrPackageNotInstalled}
			return
		}
		progress.ProgressMsg <- fmt.Sprintf("Installing package '%s'", o.PkgInstallName)
		p.InstallPackage(o, progress, tkgpackagedatamodel.OperationTypeUpdate)
		return
	}

	pkgInstallToUpdate := pkgInstall.DeepCopy()

	if pkgInstallToUpdate.Spec.PackageRef == nil || pkgInstallToUpdate.Spec.PackageRef.VersionSelection == nil {
		err = errors.New(fmt.Sprintf("failed to update package '%s' as no existing package reference/version was found in the package install", o.PkgInstallName))
		return
	}

	// If o.PackageName is provided by the user (via --package-name flag), set the package name in PackageInstall to it.
	// This is useful for the case in which the user made a typo in the package-name at the time of installation and it failed and they want to fix it through package update.
	// Otherwise if o.PackageName is not provided, fill it from the installed package spec, as the validation logic in GetPackage() needs this field to be set.
	if o.PackageName != "" {
		pkgInstallToUpdate.Spec.PackageRef.RefName = o.PackageName
	} else {
		o.PackageName = pkgInstallToUpdate.Spec.PackageRef.RefName
	}

	// If o.Version is provided by the user (via --version flag), set the version in PackageInstall to this version
	// Otherwise if o.Version is not provided, fill it from the installed package spec, as the validation logic in GetPackage() needs this field to be set.
	if o.Version != "" {
		pkgInstallToUpdate.Spec.PackageRef.VersionSelection.Constraints = o.Version
	} else {
		o.Version = pkgInstallToUpdate.Spec.PackageRef.VersionSelection.Constraints
	}

	progress.ProgressMsg <- fmt.Sprintf("Getting package metadata for '%s'", pkgInstallToUpdate.Spec.PackageRef.RefName)
	if _, _, err = p.GetPackage(o); err != nil {
		return
	}

	if secretCreated, err = p.createOrUpdateValuesSecret(o, pkgInstallToUpdate, progress.ProgressMsg); err != nil {
		return
	}

	progress.ProgressMsg <- fmt.Sprintf("Updating package install for '%s'", o.PkgInstallName)
	if err = p.kappClient.UpdatePackageInstall(pkgInstallToUpdate, secretCreated); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to update package '%s'", o.PkgInstallName))
		return
	}

	if o.Wait {
		if err = p.waitForResourceInstallation(o.PkgInstallName, o.Namespace, o.PollInterval, o.PollTimeout, progress.ProgressMsg, tkgpackagedatamodel.ResourceTypePackageInstall); err != nil {
			return
		}
	}
}

// createOrUpdateValuesSecret either creates or updates the values secret depending on whether the corresponding annotation exists or not
func (p *pkgClient) createOrUpdateValuesSecret(o *tkgpackagedatamodel.PackageOptions, pkgInstallToUpdate *kappipkg.PackageInstall, progress chan string) (bool, error) {
	var (
		secretCreated bool
		err           error
	)

	if o.ValuesFile == "" {
		return false, nil
	}

	o.SecretName = fmt.Sprintf(tkgpackagedatamodel.SecretName, o.PkgInstallName, o.Namespace)

	if o.SecretName == pkgInstallToUpdate.GetAnnotations()[tkgpackagedatamodel.TanzuPkgPluginAnnotation+"-Secret"] {
		progress <- fmt.Sprintf("Updating secret '%s'", o.SecretName)
		if err = p.updateDataValuesSecret(o); err != nil {
			err = errors.Wrap(err, "failed to update secret based on values file")
			return false, err
		}
	} else {
		progress <- fmt.Sprintf("Creating secret '%s'", o.SecretName)
		if secretCreated, err = p.createDataValuesSecret(o); err != nil {
			return secretCreated, errors.Wrap(err, "failed to create secret based on values file")
		}
	}

	pkgInstallToUpdate.Spec.Values = []kappipkg.PackageInstallValues{
		{
			SecretRef: &kappipkg.PackageInstallValuesSecretRef{Name: o.SecretName}},
	}

	return secretCreated, nil
}

// updateDataValuesSecret update a secret object containing the user-provided configuration.
func (p *pkgClient) updateDataValuesSecret(o *tkgpackagedatamodel.PackageOptions) error {
	var err error
	dataValues := make(map[string][]byte)

	if dataValues[filepath.Base(o.ValuesFile)], err = ioutil.ReadFile(o.ValuesFile); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to read from data values file '%s'", o.ValuesFile))
	}
	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: o.SecretName, Namespace: o.Namespace}, Data: dataValues,
	}

	if err := p.kappClient.GetClient().Update(context.Background(), secret); err != nil {
		return errors.Wrap(err, "failed to update Secret resource")
	}

	return nil
}
