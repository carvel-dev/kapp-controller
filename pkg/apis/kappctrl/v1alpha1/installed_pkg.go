// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type InstalledPkg struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstalledPkgSpec   `json:"spec"`
	Status InstalledPkgStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type InstalledPkgList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []InstalledPkg `json:"items"`
}

type InstalledPkgSpec struct {
	ServiceAccountName string      `json:"serviceAccountName,omitempty"`
	Cluster            *AppCluster `json:"cluster,omitempty"`

	PkgRef *InstalledPkgPkgRef  `json:"pkgRef,omitempty"`
	Values []InstalledPkgValues `json:"values,omitempty"`

	// TODO other App CR related fields
}

type InstalledPkgPkgRef struct {
	PublicName       string                           `json:"publicName,omitempty"`
	Version          string                           `json:"version,omitempty"`
	VersionSelection *versions.VersionSelectionSemver `json:"versionSelection,omitempty"`
}

type InstalledPkgValues struct {
	SecretRef *InstalledPkgValuesSecretRef `json:"secretRef,omitempty"`
}

type InstalledPkgValuesSecretRef struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

type InstalledPkgStatus struct {
	ObservedGeneration int64          `json:"observedGeneration"`
	Conditions         []AppCondition `json:"conditions"`
}
