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
type Package struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageSpec   `json:"spec"`
	Status PackageStatus `json:"status"`
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

type PackageSpec struct {
	PublicName string `json:"publicName,omitempty"`
	Version    string `json:"version,omitempty"`

	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`

	Template *v1alpha1.App `json:"template,omitempty"`
	// TODO ValuesSchema
}

type PackageStatus struct {
	ObservedGeneration int64                   `json:"observedGeneration"`
	Conditions         []v1alpha1.AppCondition `json:"conditions"`
}
