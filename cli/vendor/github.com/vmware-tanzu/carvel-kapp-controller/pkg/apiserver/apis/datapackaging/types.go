// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package datapackaging

import (
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// Package Metadata are attributes of a single package that do not change frequently and that are shared across multiple versions of a single package.
// It contains information similar to a project’s README.md.
type PackageMetadata struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PackageMetadataSpec `json:"spec"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// A package is a combination of configuration metadata and OCI images that informs the package manager what software it holds and how to install itself onto a Kubernetes cluster.
// For example, an nginx-ingress package would instruct the package manager where to download the nginx container image, how to configure the associated Deployment, and install it into a cluster.
type Package struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PackageSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PackageList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Package `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PackageMetadataList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PackageMetadata `json:"items"`
}

type PackageSpec struct {
	// The name of the PackageMetadata associated with this version
	// Must be a valid PackageMetadata name (see PackageMetadata CR for details)
	// Cannot be empty
	RefName string `json:"refName,omitempty"`
	// Package version; Referenced by PackageInstall;
	// Must be valid semver (required)
	// Cannot be empty
	Version string `json:"version,omitempty"`
	// Description of the licenses that apply to the package software
	// (optional; Array of strings)
	Licenses []string `json:"licenses,omitempty"`
	// Timestamp of release (iso8601 formatted string; optional)
	// +optional
	// +nullable
	ReleasedAt metav1.Time `json:"releasedAt,omitempty"`
	// System requirements needed to install the package.
	// Note: these requirements will not be verified by kapp-controller on
	// installation. (optional; string)
	CapactiyRequirementsDescription string `json:"capacityRequirementsDescription,omitempty"`
	// Version release notes (optional; string)
	ReleaseNotes string `json:"releaseNotes,omitempty"`

	Template AppTemplateSpec `json:"template,omitempty"`

	// valuesSchema can be used to show template values that
	// can be configured by users when a Package is installed
	// in an OpenAPI schema format.
	// +optional
	ValuesSchema ValuesSchema `json:"valuesSchema,omitempty"`

	// IncludedSoftware can be used to show the software contents of a Package.
	// This is especially useful if the underlying versions do not match the Package version
	// +optional
	IncludedSoftware []IncludedSoftware `json:"includedSoftware,omitempty"`
}

type PackageMetadataSpec struct {
	// Human friendly name of the package (optional; string)
	DisplayName string `json:"displayName,omitempty"`
	// Long description of the package (optional; string)
	LongDescription string `json:"longDescription,omitempty"`
	// Short desription of the package (optional; string)
	ShortDescription string `json:"shortDescription,omitempty"`
	// Base64 encoded icon (optional; string)
	IconSVGBase64 string `json:"iconSVGBase64,omitempty"`
	// Name of the entity distributing the package (optional; string)
	ProviderName string `json:"providerName,omitempty"`
	// List of maintainer info for the package.
	// Currently only supports the name key. (optional; array of maintner info)
	Maintainers []Maintainer `json:"maintainers,omitempty"`
	// Classifiers of the package (optional; Array of strings)
	Categories []string `json:"categories,omitempty"`
	// Description of the support available for the package (optional; string)
	SupportDescription string `json:"supportDescription,omitempty"`
}

type Maintainer struct {
	Name string `json:"name,omitempty"`
}

type AppTemplateSpec struct {
	Spec *kcv1alpha1.AppSpec `json:"spec"`
}

type ValuesSchema struct {
	// +optional
	// +nullable
	// +kubebuilder:pruning:PreserveUnknownFields
	OpenAPIv3 runtime.RawExtension `json:"openAPIv3,omitempty"`
}

// IncludedSoftware contains the underlying Software Contents of a Package
type IncludedSoftware struct {
	DisplayName string `json:"displayName,omitempty" protobuf:"bytes,1,opt,name=displayName"`
	Version     string `json:"version,omitempty" protobuf:"bytes,2,opt,name=version"`
	Description string `json:"description,omitempty" protobuf:"bytes,3,opt,name=description"`
}
