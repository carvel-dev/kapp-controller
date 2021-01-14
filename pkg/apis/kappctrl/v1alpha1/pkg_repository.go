// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PkgRepository struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PkgRepositorySpec   `json:"spec"`
	Status PkgRepositoryStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PkgRepositoryList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PkgRepository `json:"items"`
}

type PkgRepositorySpec struct {
	Fetch *PkgRepositoryFetch `json:"fetch"`
}

type PkgRepositoryFetch struct {
	Image *AppFetchImage `json:"image,omitempty"`
	HTTP  *AppFetchHTTP  `json:"http,omitempty"`
	Git   *AppFetchGit   `json:"git,omitempty"`
}

type PkgRepositoryStatus struct {
	ObservedGeneration int64          `json:"observedGeneration"`
	Conditions         []AppCondition `json:"conditions"`
}
