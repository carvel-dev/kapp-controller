// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"

	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
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
func parseRegistryImageURL(imgURL string) (repository, tag string, err error) {
	ref, err := name.ParseReference(imgURL, name.WeakValidation)
	if err != nil {
		return "", "", err
	}

	repository = ref.Context().String()
	tag = ref.Identifier()
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
