// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"
	kappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	defaultRepositoryImageTag = "latest"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "repository",
		Aliases: []string{"pkgr", "repo"},
		Short:   "Repository",
	}
	return cmd
}

// parseRegistryImageURL parses the registry image URL to get repository and tag, tag is empty if not specified
func parseRegistryImageURL(imgURL string) (string, string, error) {
	ref, err := name.ParseReference(imgURL, name.WeakValidation)
	if err != nil {
		return "", "", err
	}

	repository := ref.Context().String()
	tag := ref.Identifier()

	// the parser function sets the tag to "latest" if not specified, however we want it to be empty
	if tag == defaultRepositoryImageTag && !strings.HasSuffix(imgURL, ":"+defaultRepositoryImageTag) {
		tag = ""
	}
	return repository, tag, nil
}

// getCurrentRepositoryAndTagInUse fetches the current tag used by package repository, taking tagselection into account
// TODO: Should we error out if the fetch does not pull an imgpkg bundle?
func getCurrentRepositoryAndTagInUse(pkgr *kappipkg.PackageRepository) (repository, tag string, err error) {
	if pkgr.Spec.Fetch == nil || pkgr.Spec.Fetch.ImgpkgBundle == nil {
		return "", "", fmt.Errorf("Failed to find OCI registry URL")
	}

	repository, tag, err = parseRegistryImageURL(pkgr.Spec.Fetch.ImgpkgBundle.Image)
	if err != nil {
		return "", "", fmt.Errorf("Failed to parse OCI registry URL: %s", err.Error())
	}

	if tag == "" {
		tag = defaultRepositoryImageTag
	}

	if pkgr.Spec.Fetch.ImgpkgBundle.TagSelection != nil && pkgr.Spec.Fetch.ImgpkgBundle.TagSelection.Semver != nil {
		tag = fmt.Sprintf("(%s)", pkgr.Spec.Fetch.ImgpkgBundle.TagSelection.Semver.Constraints)
	}

	return repository, tag, nil
}

func newPackageRepository(name, url, namespace string) (*v1alpha1.PackageRepository, error) {
	pkgr := &v1alpha1.PackageRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	return updateExistingPackageRepository(pkgr, url)
}

func updateExistingPackageRepository(pkgr *v1alpha1.PackageRepository,
	url string) (*v1alpha1.PackageRepository, error) {

	pkgr = pkgr.DeepCopy()

	pkgr.Spec = kappipkg.PackageRepositorySpec{
		Fetch: &kappipkg.PackageRepositoryFetch{
			ImgpkgBundle: &kappctrl.AppFetchImgpkgBundle{Image: url},
		},
	}

	_, tag, err := parseRegistryImageURL(url)
	if err != nil {
		return nil, fmt.Errorf("Parsing OCI registry URL: %s", err)
	}

	if tag == "" {
		pkgr.Spec.Fetch.ImgpkgBundle.TagSelection = &versions.VersionSelection{
			Semver: &versions.VersionSelectionSemver{},
		}
	}

	return pkgr, err
}

func waitForPackageRepositoryInstallation(pollInterval time.Duration, pollTimeout time.Duration,
	namespace string, repository string, client versioned.Interface) error {
	var (
		status             kappctrl.GenericStatus
		reconcileSucceeded bool
	)
	if err := wait.Poll(pollInterval, pollTimeout, func() (done bool, err error) {
		resource, err := client.PackagingV1alpha1().PackageRepositories(
			namespace).Get(context.Background(), repository, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if resource.Generation != resource.Status.ObservedGeneration {
			// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
			return false, nil
		}
		status = resource.Status.GenericStatus

		for _, cond := range status.Conditions {
			switch {
			case cond.Type == kappctrl.ReconcileSucceeded && cond.Status == corev1.ConditionTrue:
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
		return fmt.Errorf("PackageRepository reconciliation failed")
	}
	return nil
}
