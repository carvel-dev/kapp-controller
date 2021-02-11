// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/package/runtime.Object
type InstalledPackage struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstalledPackageSpec   `json:"spec"`
	Status InstalledPackageStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/package/runtime.Object
type InstalledPackageList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []InstalledPackage `json:"items"`
}

type InstalledPackageSpec struct {
	ServiceAccountName string               `json:"serviceAccountName,omitempty"`
	Cluster            *v1alpha1.AppCluster `json:"cluster,omitempty"`

	PkgRef *PackageRef              `json:"pkgRef,omitempty"`
	Values []InstalledPackageValues `json:"values,omitempty"`

	// TODO other App CR related fields
}

type PackageRef struct {
	PublicName       string                           `json:"publicName,omitempty"`
	Version          string                           `json:"version,omitempty"`
	VersionSelection *versions.VersionSelectionSemver `json:"versionSelection,omitempty"`
}

type InstalledPackageValues struct {
	SecretRef *InstalledPackageValuesSecretRef `json:"secretRef,omitempty"`
}

type InstalledPackageValuesSecretRef struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

type InstalledPackageStatus struct {
	v1alpha1.GenericStatus `json:",inline"`
	// TODO this is desired resolved version (not actually deployed)
	Version string `json:"version,omitempty"`
}
