// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"fmt"
	"strings"

	"carvel.dev/vendir/pkg/vendir/versions"
	semver "github.com/k14s/semver/v4"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// package validations
func ValidatePackageMetadata(pkgm datapackaging.PackageMetadata) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidatePackageMetadataName(pkgm.ObjectMeta.Name, field.NewPath("metadata").Child("name"))...)

	return allErrs
}

// validate name
func ValidatePackageMetadataName(pkgmName string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs,
		validation.IsFullyQualifiedName(fldPath, pkgmName)...)

	return allErrs
}

// package version validations

func ValidatePackage(pkg datapackaging.Package) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs,
		ValidatePackageSpecPackageName(pkg.Spec.RefName, field.NewPath("spec", "refName"))...)

	allErrs = append(allErrs, ValidatePackageSpecVersion(pkg.Spec.Version, field.NewPath("spec", "version"))...)

	allErrs = append(allErrs,
		ValidatePackageName(pkg.ObjectMeta.Name, pkg.Spec.RefName, pkg.Spec.Version, field.NewPath("metadata").Child("name"))...)

	if pkg.Spec.KubernetesVersionSelection != nil {
		allErrs = append(allErrs,
			ValidatePackageVersionConstraints(pkg.Spec.KubernetesVersionSelection.Constraints, field.NewPath("spec", "kubernetesVersionSelection", "constraints"))...)
	}

	if pkg.Spec.KappControllerVersionSelection != nil {
		allErrs = append(allErrs,
			ValidatePackageVersionConstraints(pkg.Spec.KappControllerVersionSelection.Constraints, field.NewPath("spec", "kappControllerVersionSelection", "constraints"))...)
	}

	return allErrs
}

// ValidatePackageVersionConstraints checks that the user given constrain range is valid
func ValidatePackageVersionConstraints(constraints string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	_, err := semver.ParseRange(constraints)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, constraints, err.Error()))
	}

	return allErrs
}

// validate metdata.name = spec.RefName + spec.Version
func ValidatePackageName(pkgName, pkgmName, pkgVersion string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if !(pkgName == pkgmName+"."+pkgVersion) {
		allErrs = append(allErrs,
			field.Invalid(fldPath, pkgName, "must be <spec.refName> + '.' + <spec.version>"))
	}

	return allErrs
}

// validate spec.version is not empty and is a valid semver
func ValidatePackageSpecVersion(version string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if version == "" {
		allErrs = append(allErrs,
			field.Invalid(fldPath, version, "cannot be empty"))
	}

	if _, err := versions.NewSemver(version); err != nil {
		allErrs = append(allErrs,
			field.Invalid(fldPath, version, fmt.Sprintf("must be valid semver: %v", err)))
	}

	return allErrs
}

// validate spec.RefName isnt empty
func ValidatePackageSpecPackageName(name string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if name == "" {
		allErrs = append(allErrs, field.Required(fldPath, "can not be empty"))
	}

	allErrs = append(allErrs, ValidatePackageMetadataName(name, fldPath)...)
	return allErrs
}

func IsFullyQualifiedName(fldPath *field.Path, name string) field.ErrorList {
	var allErrors field.ErrorList
	if len(name) == 0 {
		return append(allErrors, field.Required(fldPath, ""))
	}
	if errs := validation.IsDNS1123Subdomain(name); len(errs) > 0 {
		return append(allErrors, field.Invalid(fldPath, name, strings.Join(errs, ",")))
	}
	if len(strings.Split(name, ".")) < 3 {
		return append(allErrors, field.Invalid(fldPath, name, "should be a domain with at least three segments separated by dots"))
	}
	return allErrors
}
