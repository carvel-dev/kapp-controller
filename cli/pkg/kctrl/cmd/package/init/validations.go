// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"
	"strings"

	"github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions"
	"k8s.io/apimachinery/pkg/util/validation"
)

//TODO should we use the same validation used in kapp controller. But that accepts other parameter. ValidatePackageMetadataName in validations.go file
func validateFQName(name string) (bool, string, error) {
	if len(name) == 0 {
		return false, "Fully qualified name of a package cannot be empty", nil
	}
	if errs := validation.IsDNS1123Subdomain(name); len(errs) > 0 {
		return false, strings.Join(errs, ","), nil
	}
	if len(strings.Split(name, ".")) < 3 {
		return false, fmt.Sprintf("Invalid name: %s should be a fully qualified name with at least three segments separated by dots", name), nil
	}
	return true, "", nil
}

func validatePackageSpecVersion(version string) (bool, string, error) {
	if version == "" {
		return false, "Version cannot be empty", nil
	}
	if _, err := versions.NewSemver(version); err != nil {
		return false, fmt.Sprintf("Invalid version: %s must be valid semver: %v", version, err), nil
	}
	return true, "", nil
}
