// Copyright 2021 VMware, Inc.
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
	RefName  string   `json:"refName,omitempty" protobuf:"bytes,1,opt,name=refName"`
	Version  string   `json:"version,omitempty" protobuf:"bytes,2,opt,name=version"`
	Licenses []string `json:"licenses,omitempty" protobuf:"bytes,3,rep,name=licenses"`
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
