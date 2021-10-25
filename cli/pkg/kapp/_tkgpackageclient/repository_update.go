// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"fmt"

	kappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"

	"github.com/pkg/errors"
	k8serror "k8s.io/apimachinery/pkg/api/errors"

	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

func (p *pkgClient) UpdateRepository(o *tkgpackagedatamodel.RepositoryOptions, progress *tkgpackagedatamodel.PackageProgress) {
	var (
		existingRepository *kappipkg.PackageRepository
		err                error
		tag                string
	)

	defer func() {
		progressCleanup(err, progress)
	}()

	progress.ProgressMsg <- fmt.Sprintf("Getting package repository '%s'", o.RepositoryName)
	existingRepository, err = p.kappClient.GetPackageRepository(o.RepositoryName, o.Namespace)
	if err != nil {
		if k8serror.IsNotFound(err) {
			err = nil
		} else {
			return
		}
	}

	if existingRepository != nil {
		repositoryToUpdate := existingRepository.DeepCopy()
		if err = p.validateRepositoryUpdate(o.RepositoryName, o.RepositoryURL, o.Namespace); err != nil {
			return
		}

		_, tag, err = parseRegistryImageURL(o.RepositoryURL)
		if err != nil {
			err = errors.Wrap(err, "failed to parse OCI registry URL")
			return
		}

		repositoryToUpdate.Spec = kappipkg.PackageRepositorySpec{
			Fetch: &kappipkg.PackageRepositoryFetch{
				ImgpkgBundle: &kappctrl.AppFetchImgpkgBundle{Image: o.RepositoryURL},
			},
		}

		if tag == "" {
			repositoryToUpdate.Spec.Fetch.ImgpkgBundle.TagSelection = &versions.VersionSelection{
				Semver: &versions.VersionSelectionSemver{
					Constraints: tkgpackagedatamodel.DefaultRepositoryImageTagConstraint,
				},
			}
		}

		progress.ProgressMsg <- "Updating package repository resource"
		if err = p.kappClient.UpdatePackageRepository(repositoryToUpdate); err != nil {
			err = errors.Wrap(err, fmt.Sprintf("failed to update package repository '%s' in namespace '%s'", o.RepositoryName, o.Namespace))
			return
		}

		if o.Wait {
			if err = p.waitForResourceInstallation(o.RepositoryName, o.Namespace, o.PollInterval, o.PollTimeout, progress.ProgressMsg, tkgpackagedatamodel.ResourceTypePackageRepository); err != nil {
				return
			}
		}
	} else if o.CreateRepository {
		p.AddRepository(o, progress, tkgpackagedatamodel.OperationTypeUpdate)
	} else {
		err = &tkgpackagedatamodel.PackagePluginNonCriticalError{Reason: tkgpackagedatamodel.ErrRepoNotExists}
	}
}

// validateRepositoryUpdate ensures that another repository (with the same OCI registry URL) does not already exist in the cluster
func (p *pkgClient) validateRepositoryUpdate(repositoryName, repositoryImg, namespace string) error {
	repositoryList, err := p.kappClient.ListPackageRepositories(namespace)
	if err != nil {
		return errors.Wrap(err, "failed to list package repositories")
	}

	for _, repository := range repositoryList.Items { //nolint:gocritic
		// This stops the update validation to compare with itself
		if repository.Name == repositoryName {
			continue
		}

		if repository.Spec.Fetch != nil && repository.Spec.Fetch.ImgpkgBundle != nil &&
			repository.Spec.Fetch.ImgpkgBundle.Image == repositoryImg {
			return errors.New("repository with the same OCI registry URL already exists")
		}
	}

	return nil
}
