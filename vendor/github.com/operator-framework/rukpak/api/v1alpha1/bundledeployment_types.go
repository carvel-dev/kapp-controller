/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	BundleDeploymentGVK  = SchemeBuilder.GroupVersion.WithKind("BundleDeployment")
	BundleDeploymentKind = BundleDeploymentGVK.Kind
)

const (
	TypeHasValidBundle = "HasValidBundle"
	TypeInstalled      = "Installed"

	ReasonBundleLoadFailed         = "BundleLoadFailed"
	ReasonReadingContentFailed     = "ReadingContentFailed"
	ReasonErrorGettingClient       = "ErrorGettingClient"
	ReasonErrorGettingReleaseState = "ErrorGettingReleaseState"
	ReasonInstallFailed            = "InstallFailed"
	ReasonUpgradeFailed            = "UpgradeFailed"
	ReasonReconcileFailed          = "ReconcileFailed"
	ReasonCreateDynamicWatchFailed = "CreateDynamicWatchFailed"
	ReasonInstallationSucceeded    = "InstallationSucceeded"
)

// BundleDeploymentSpec defines the desired state of BundleDeployment
type BundleDeploymentSpec struct {
	//+kubebuilder:validation:Pattern:=^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
	// ProvisionerClassName sets the name of the provisioner that should reconcile this BundleDeployment.
	ProvisionerClassName string `json:"provisionerClassName"`
	// Template describes the generated Bundle that this deployment will manage.
	Template *BundleTemplate `json:"template"`
	// Config is provisioner specific configurations
	// +kubebuilder:pruning:PreserveUnknownFields
	Config runtime.RawExtension `json:"config,omitempty"`
}

// BundleTemplate defines the desired state of a Bundle resource
type BundleTemplate struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the desired behavior of the Bundle.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec BundleSpec `json:"spec"`
}

// BundleDeploymentStatus defines the observed state of BundleDeployment
type BundleDeploymentStatus struct {
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
	ActiveBundle       string             `json:"activeBundle,omitempty"`
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,shortName={"bd","bds"}
//+kubebuilder:printcolumn:name="Active Bundle",type=string,JSONPath=`.status.activeBundle`
//+kubebuilder:printcolumn:name="Install State",type=string,JSONPath=`.status.conditions[?(.type=="Installed")].reason`
//+kubebuilder:printcolumn:name=Age,type=date,JSONPath=`.metadata.creationTimestamp`
//+kubebuilder:printcolumn:name=Provisioner,type=string,JSONPath=`.spec.provisionerClassName`,priority=1

// BundleDeployment is the Schema for the bundledeployments API
type BundleDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BundleDeploymentSpec   `json:"spec"`
	Status BundleDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BundleDeploymentList contains a list of BundleDeployment
type BundleDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BundleDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BundleDeployment{}, &BundleDeploymentList{})
}
