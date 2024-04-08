// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PackageMetadata struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec PackageMetadataSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Package struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec PackageSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PackageList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []Package `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PackageMetadataList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []PackageMetadata `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type PackageSpec struct {
	RefName string `json:"refName,omitempty" protobuf:"bytes,1,opt,name=refName"`
	Version string `json:"version,omitempty" protobuf:"bytes,2,opt,name=version"`
	// List of dependencies to be resolved
	Dependencies []DependencyType `json:"dependencies,omitempty" protobuf:"bytes,12,rep,name=dependencies"`
	Licenses     []string         `json:"licenses,omitempty" protobuf:"bytes,3,rep,name=licenses"`
	// +optional
	// +nullable
	ReleasedAt                      metav1.Time `json:"releasedAt,omitempty" protobuf:"bytes,4,opt,name=releasedAt"`
	CapactiyRequirementsDescription string      `json:"capacityRequirementsDescription,omitempty" protobuf:"bytes,5,opt,name=capacityRequirementsDescription"`
	ReleaseNotes                    string      `json:"releaseNotes,omitempty" protobuf:"bytes,6,opt,name=releaseNotes"`

	Template AppTemplateSpec `json:"template,omitempty" protobuf:"bytes,7,opt,name=template"`

	// valuesSchema can be used to show template values that
	// can be configured by users when a Package is installed
	// in an OpenAPI schema format.
	// +optional
	ValuesSchema ValuesSchema `json:"valuesSchema,omitempty" protobuf:"bytes,8,opt,name=valuesSchema"`

	// IncludedSoftware can be used to show the software contents of a Package.
	// This is especially useful if the underlying versions do not match the Package version
	// +optional
	IncludedSoftware []IncludedSoftware `json:"includedSoftware,omitempty" protobuf:"bytes,9,opt,name=includedSoftware"`

	// KappControllerVersionSelection specifies the versions of kapp-controller which can install this package
	// +optional
	KappControllerVersionSelection *VersionSelection `json:"kappControllerVersionSelection,omitempty" protobuf:"bytes,10,opt,name=kappControllerVersionSelection"`
	// KubernetesVersionSelection specifies the versions of k8s which this package can be installed on
	// +optional
	KubernetesVersionSelection *VersionSelection `json:"kubernetesVersionSelection,omitempty" protobuf:"bytes,11,opt,name=kubernetesVersionSelection"`
}

// VersionSelection provides version range constraints but will always accept prereleases
type VersionSelection struct {
	Constraints string `json:"constraints,omitempty" protobuf:"bytes,1,opt,name=constraints"`
}

type PackageMetadataSpec struct {
	DisplayName        string       `json:"displayName,omitempty" protobuf:"bytes,1,opt,name=displayName"`
	LongDescription    string       `json:"longDescription,omitempty" protobuf:"bytes,2,opt,name=longDescription"`
	ShortDescription   string       `json:"shortDescription,omitempty" protobuf:"bytes,3,opt,name=shortDescription"`
	IconSVGBase64      string       `json:"iconSVGBase64,omitempty" protobuf:"bytes,4,opt,name=iconSVGBase64"`
	ProviderName       string       `json:"providerName,omitempty" protobuf:"bytes,5,opt,name=providerName"`
	Maintainers        []Maintainer `json:"maintainers,omitempty" protobuf:"bytes,6,rep,name=maintainers"`
	Categories         []string     `json:"categories,omitempty" protobuf:"bytes,7,rep,name=categories"`
	SupportDescription string       `json:"supportDescription,omitempty" protobuf:"bytes,8,opt,name=supportDescription"`
}

type Maintainer struct {
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
}

type AppTemplateSpec struct {
	Spec *kcv1alpha1.AppSpec `json:"spec" protobuf:"bytes,1,opt,name=spec"`
}

type ValuesSchema struct {
	// +optional
	// +nullable
	// +kubebuilder:pruning:PreserveUnknownFields
	OpenAPIv3 runtime.RawExtension `json:"openAPIv3,omitempty" protobuf:"bytes,1,opt,name=openAPIv3"`
}

// IncludedSoftware contains the underlying Software Contents of a Package
type IncludedSoftware struct {
	DisplayName string `json:"displayName,omitempty" protobuf:"bytes,1,opt,name=displayName"`
	Version     string `json:"version,omitempty" protobuf:"bytes,2,opt,name=version"`
	Description string `json:"description,omitempty" protobuf:"bytes,3,opt,name=description"`
}

// DependencyType contains the list of type of depencies
type DependencyType struct {
	Package *PackageDependency `json:"package,omitempty" protobuf:"bytes,1,opt,name=package"`
}

// PackageDependency contains package dependency related info
type PackageDependency struct {
	// The name of the PackageMetadata associated with this dependency
	// Must be a valid PackageMetadata name (see PackageMetadata CR for details)
	// Cannot be empty
	RefName string `json:"refName,omitempty" protobuf:"bytes,1,opt,name=refName"`
	// Package version; Will be Referenced by PackageInstall;
	// Must be valid semver (required)
	// Cannot be empty
	Version string `json:"version,omitempty" protobuf:"bytes,2,opt,name=version"`
}
