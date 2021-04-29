// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name=Description,JSONPath=.status.friendlyDescription,description=Friendly description,type=string
// +kubebuilder:printcolumn:name=Since-Deploy,JSONPath=.status.deploy.startedAt,description=Last time app started being deployed. Does not mean anything was changed.,type=date
// +kubebuilder:printcolumn:name=Age,JSONPath=.metadata.creationTimestamp,description=Time since creation,type=date
type App struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +optional
	Spec AppSpec `json:"spec,omitempty"`
	// +optional
	Status AppStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AppList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []App `json:"items"`
}

// +k8s:openapi-gen=true
type AppSpec struct {
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// +optional
	Cluster *AppCluster `json:"cluster,omitempty"`
	// +optional
	Fetch []AppFetch `json:"fetch,omitempty"`
	// +optional
	Template []AppTemplate `json:"template,omitempty"`
	// +optional
	Deploy []AppDeploy `json:"deploy,omitempty"`
	// Paused when set to true will ignore all pending changes,
	// once it set back to false, pending changes will be applied
	// +optional
	Paused bool `json:"paused,omitempty"`
	// Canceled when set to true will stop all active changes
	// +optional
	Canceled bool `json:"canceled,omitempty"`
	// Controls frequency of app reconciliation
	// +optional
	SyncPeriod *metav1.Duration `json:"syncPeriod,omitempty"`
	// When NoopDeletion set to true, App deletion should
	// delete App CR but preserve App's associated resources
	// +optional
	NoopDelete bool `json:"noopDelete,omitempty"`
}

// +k8s:openapi-gen=true
type AppCluster struct {
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// +optional
	KubeconfigSecretRef *AppClusterKubeconfigSecretRef `json:"kubeconfigSecretRef,omitempty"`
}

// +k8s:openapi-gen=true
type AppClusterKubeconfigSecretRef struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

type AppStatus struct {
	// +optional
	ManagedAppName string `json:"managedAppName,omitempty"`
	// +optional
	Fetch *AppStatusFetch `json:"fetch,omitempty"`
	// +optional
	Template *AppStatusTemplate `json:"template,omitempty"`
	// +optional
	Deploy *AppStatusDeploy `json:"deploy,omitempty"`
	// +optional
	Inspect *AppStatusInspect `json:"inspect,omitempty"`
	// +optional
	ConsecutiveReconcileSuccesses int `json:"consecutiveReconcileSuccesses,omitempty"`
	// +optional
	ConsecutiveReconcileFailures int `json:"consecutiveReconcileFailures,omitempty"`

	GenericStatus `json:",inline"`
}

type AppStatusFetch struct {
	// +optional
	Stderr string `json:"stderr,omitempty"`
	// +optional
	Stdout string `json:"stdout,omitempty"`
	// +optional
	ExitCode int `json:"exitCode"`
	// +optional
	Error string `json:"error,omitempty"`
	// +optional
	StartedAt metav1.Time `json:"startedAt,omitempty"`
	// +optional
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}

type AppStatusTemplate struct {
	// +optional
	Stderr string `json:"stderr,omitempty"`
	// +optional
	ExitCode int `json:"exitCode"`
	// +optional
	Error string `json:"error,omitempty"`
	// +optional
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}

type AppStatusDeploy struct {
	// +optional
	Stdout string `json:"stdout,omitempty"`
	// +optional
	Stderr string `json:"stderr,omitempty"`
	// +optional
	Finished bool `json:"finished"`
	// +optional
	ExitCode int `json:"exitCode"`
	// +optional
	Error string `json:"error,omitempty"`
	// +optional
	StartedAt metav1.Time `json:"startedAt,omitempty"`
	// +optional
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}

type AppStatusInspect struct {
	// +optional
	Stdout string `json:"stdout,omitempty"`
	// +optional
	Stderr string `json:"stderr,omitempty"`
	// +optional
	ExitCode int `json:"exitCode"`
	// +optional
	Error string `json:"error,omitempty"`
	// +optional
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}
