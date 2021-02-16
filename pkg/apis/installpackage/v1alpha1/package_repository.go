// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PackageRepository struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageRepositorySpec   `json:"spec"`
	Status PackageRepositoryStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PackageRepositoryList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PackageRepository `json:"items"`
}

type PackageRepositorySpec struct {
	Fetch *PackageRepositoryFetch `json:"fetch"`
}

type PackageRepositoryFetch struct {
	Image  *v1alpha1.AppFetchImage        `json:"image,omitempty"`
	HTTP   *v1alpha1.AppFetchHTTP         `json:"http,omitempty"`
	Git    *v1alpha1.AppFetchGit          `json:"git,omitempty"`
	Bundle *v1alpha1.AppFetchImgpkgBundle `json:"bundle,omitempty"`
}

type PackageRepositoryStatus struct {
	ObservedGeneration int64                   `json:"observedGeneration"`
	Conditions         []v1alpha1.AppCondition `json:"conditions"`
}
