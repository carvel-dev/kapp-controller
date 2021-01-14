// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Pkg struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PkgSpec   `json:"spec"`
	Status PkgStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PkgList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Pkg `json:"items"`
}

type PkgSpec struct {
	PublicName string `json:"publicName,omitempty"`
	Version    string `json:"version,omitempty"`

	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`

	Template *App `json:"template,omitempty"`
	// TODO ValuesSchema
}

type PkgStatus struct {
	ObservedGeneration int64          `json:"observedGeneration"`
	Conditions         []AppCondition `json:"conditions"`
}
