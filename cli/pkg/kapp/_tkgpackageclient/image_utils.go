// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/pkg/errors"

	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

// parseRegistryImageURL parses the registry image URL to get repository and tag, tag is empty if not specified
func parseRegistryImageURL(imgURL string) (repository, tag string, err error) {
	ref, err := name.ParseReference(imgURL, name.WeakValidation)
	if err != nil {
		return "", "", err
	}

	repository = ref.Context().String()
	tag = ref.Identifier()
	// the parser function sets the tag to "latest" if not specified, however we want it to be empty
	if tag == tkgpackagedatamodel.DefaultRepositoryImageTag && !strings.HasSuffix(imgURL, ":"+tkgpackagedatamodel.DefaultRepositoryImageTag) {
		tag = ""
	}
	return repository, tag, nil
}

// GetCurrentRepositoryAndTagInUse fetches the current tag used by package repository, taking tagselection into account
func GetCurrentRepositoryAndTagInUse(pkgr *kappipkg.PackageRepository) (repository, tag string, err error) {
	if pkgr.Spec.Fetch == nil || pkgr.Spec.Fetch.ImgpkgBundle == nil {
		return "", "", errors.New("failed to find OCI registry URL")
	}

	repository, tag, err = parseRegistryImageURL(pkgr.Spec.Fetch.ImgpkgBundle.Image)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to parse OCI registry URL")
	}

	if tag == "" {
		tag = tkgpackagedatamodel.DefaultRepositoryImageTag
	}

	if pkgr.Spec.Fetch.ImgpkgBundle.TagSelection != nil && pkgr.Spec.Fetch.ImgpkgBundle.TagSelection.Semver != nil {
		tag = fmt.Sprintf("(%s)", pkgr.Spec.Fetch.ImgpkgBundle.TagSelection.Semver.Constraints)
	}

	return repository, tag, nil
}
