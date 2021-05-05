// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"fmt"
	"regexp"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/packages"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type InternalPackage packages.Package

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type InternalPackageList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []InternalPackage `json:"items"`
}

func (ip InternalPackage) ValidatePackageName() error {
	const pkgNameRegex = ".*\\..*\\..*"
	reg := regexp.MustCompile(pkgNameRegex)
	if !reg.MatchString(ip.Name) {
		return fmt.Errorf("package name requires at least two periods. Recommended package naming convention is publicName.packageRepo.version")
	}
	return nil
}
