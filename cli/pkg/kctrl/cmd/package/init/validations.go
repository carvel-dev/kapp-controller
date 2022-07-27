// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"
)

// TODO should we use the same validation used in kapp controller. But that accepts other parameter. ValidatePackageMetadataName in validations.go file
func validateRefName(name string) (bool, string, error) {
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
