// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"fmt"

	"github.com/pkg/errors"

	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	kapppkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

// GetPackageInstall takes an installed package name and returns the corresponding installed package
func (p *pkgClient) GetPackageInstall(o *tkgpackagedatamodel.PackageOptions) (*kappipkg.PackageInstall, error) {
	pkg, err := p.kappClient.GetPackageInstall(o.PackageName, o.Namespace)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to find installed package '%s' in namespace '%s'", o.PackageName, o.Namespace))
	}
	return pkg, nil
}

// GetPackage takes a package name and package version and returns the corresponding PackageMetadata and Package.
// If the resolution is unsuccessful, an error is returned.
func (p *pkgClient) GetPackage(o *tkgpackagedatamodel.PackageOptions) (*kapppkg.PackageMetadata, *kapppkg.Package, error) {
	var (
		resolvedPackage *kapppkg.PackageMetadata
		err             error
	)

	if resolvedPackage, err = p.kappClient.GetPackageMetadataByName(o.PackageName, o.Namespace); err != nil {
		return nil, nil, errors.Wrapf(err, "failed to find a package with name '%s' in namespace '%s'", o.PackageName, o.Namespace)
	}

	packageVersions, err := p.kappClient.ListPackages(o.PackageName, o.Namespace)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to list package versions")
	}

	for _, item := range packageVersions.Items { //nolint:gocritic
		if item.Spec.Version == o.Version {
			return resolvedPackage, &item, nil
		}
	}

	return nil, nil, errors.Errorf(fmt.Sprintf("failed to resolve package '%s' with version '%s'", o.PackageName, o.Version))
}
