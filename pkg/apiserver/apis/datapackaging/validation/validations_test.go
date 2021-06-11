// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package validation_test

import (
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidatePackageMetadataNameInvalid(t *testing.T) {
	invalidName := "bummer-boy"
	// Name could be invalid for many reasons so just assert we have
	// an error relating to name and not specific error string
	expectedErr := field.Error{
		Type:  field.ErrorTypeInvalid,
		Field: "metadata.name",
	}

	errList := validation.ValidatePackageMetadataName(invalidName, field.NewPath("metadata").Child("name"))

	if len(errList) == 0 {
		t.Fatalf("Expected validation to error when given invalid name")
	}

	if !contains(errList, expectedErr) {
		t.Fatalf("Expected invalid field error for metadata.name, but got: %v", errList.ToAggregate())
	}
}

func TestValidatePackageMetadataNameValid(t *testing.T) {
	validName := "package.carvel.dev"

	errList := validation.ValidatePackageMetadataName(validName, field.NewPath("metadata").Child("name"))

	if len(errList) != 0 {
		t.Fatalf("Expected no error for valid name, but got: %v", errList.ToAggregate())
	}
}

func TestValidatePackageVersionNameInvalid(t *testing.T) {
	invalidName := "pkg.3.0"
	pkgName := "pkg"
	pkgVersion := "2.0"
	expectedErr := field.Error{
		Type:  field.ErrorTypeInvalid,
		Field: "metadata.name",
	}

	errList := validation.ValidatePackageName(invalidName, pkgName, pkgVersion, field.NewPath("metadata", "name"))

	if len(errList) == 0 {
		t.Fatalf("Expected error when PackageVersion name is invalid")
	}

	if !contains(errList, expectedErr) {
		t.Fatalf("Expected invalid field error for metadata.name, but got: %v", errList.ToAggregate())
	}
}

func TestValidatePackageVersionNameValid(t *testing.T) {
	validName := "pkg.2.0"
	pkgName := "pkg"
	pkgVersion := "2.0"

	errList := validation.ValidatePackageName(validName, pkgName, pkgVersion, field.NewPath("metadata", "name"))

	if len(errList) != 0 {
		t.Fatalf("Expected no error when PackageVersion name is valid, but got: %v", errList.ToAggregate())
	}
}

func TestValidatePackageVersionSpecPackageVersionInvalid(t *testing.T) {
	invalidVersion := ""
	expectedErr := field.Error{
		Type:  field.ErrorTypeInvalid,
		Field: "spec.version",
	}

	errList := validation.ValidatePackageSpecVersion(invalidVersion, field.NewPath("spec", "version"))

	if len(errList) == 0 {
		t.Fatalf("Expected error when spec.version is invalid")
	}

	if !contains(errList, expectedErr) {
		t.Fatalf("Expected invalid field error for spec.version, but got: %v", errList.ToAggregate())
	}
}

func TestValidatePackageVersionSpecPackageVersionValid(t *testing.T) {
	validVersion := "1.0.0"

	errList := validation.ValidatePackageSpecVersion(validVersion, field.NewPath("spec", "version"))

	if len(errList) != 0 {
		t.Fatalf("Expected no error when spec.version is valid")
	}
}

func TestValidatePackageVersionSpecPackageNameInvalid(t *testing.T) {
	invalidName := ""
	expectedErr := field.Error{
		Type:  field.ErrorTypeRequired,
		Field: "spec.packageName",
	}

	errList := validation.ValidatePackageSpecPackageName(invalidName, field.NewPath("spec", "packageName"))

	if len(errList) == 0 {
		t.Fatalf("Expected error when spec.packageName is invalid")
	}

	if !contains(errList, expectedErr) {
		t.Fatalf("Expected invalid field error for spec.packageName, but got: %v", errList.ToAggregate())
	}
}

func TestValidatePackageVersionSpecPackageNameValid(t *testing.T) {
	validName := "package.carvel.dev"

	errList := validation.ValidatePackageSpecPackageName(validName, field.NewPath("spec", "packageName"))

	if len(errList) != 0 {
		t.Fatalf("Expected no error when spec.packageName is valid")
	}
}

// Searchs for Error in ErrorList by Type + Field, but not details
func contains(errList field.ErrorList, expectedErr field.Error) bool {
	for _, err := range errList {
		if err.Type == expectedErr.Type && err.Field == expectedErr.Field {
			return true
		}
	}
	return false
}
