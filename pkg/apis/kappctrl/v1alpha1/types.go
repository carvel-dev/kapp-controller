// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type App struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppSpec   `json:"spec"`
	Status AppStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AppList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []App `json:"items"`
}

type AppSpec struct {
	ServiceAccountName string      `json:"serviceAccountName,omitempty"`
	Cluster            *AppCluster `json:"cluster,omitempty"`

	Fetch    []AppFetch    `json:"fetch,omitempty"`
	Template []AppTemplate `json:"template,omitempty"`
	Deploy   []AppDeploy   `json:"deploy,omitempty"`

	// Paused when set to true will ignore all pending changes,
	// once it set back to false, pending changes will be applied
	Paused bool `json:"paused,omitempty"`
	// Canceled when set to true will stop all active changes
	Canceled bool `json:"canceled,omitempty"`
	// Controls frequency of app reconciliation
	SyncPeriod *metav1.Duration `json:"syncPeriod,omitempty"`
}

type AppCluster struct {
	Namespace           string                         `json:"namespace,omitempty"`
	KubeconfigSecretRef *AppClusterKubeconfigSecretRef `json:"kubeconfigSecretRef,omitempty"`
}

type AppClusterKubeconfigSecretRef struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

type AppStatus struct {
	ManagedAppName string `json:"managedAppName,omitempty"`

	Fetch    *AppStatusFetch    `json:"fetch,omitempty"`
	Template *AppStatusTemplate `json:"template,omitempty"`
	Deploy   *AppStatusDeploy   `json:"deploy,omitempty"`
	Inspect  *AppStatusInspect  `json:"inspect,omitempty"`

	ConsecutiveReconcileSuccesses int `json:"consecutiveReconcileSuccesses,omitempty"`
	ConsecutiveReconcileFailures  int `json:"consecutiveReconcileFailures,omitempty"`

	GenericStatus `json:",inline"`
}

type AppStatusFetch struct {
	Stderr    string      `json:"stderr,omitempty"`
	Stdout    string      `json:"stdout,omitempty"`
	ExitCode  int         `json:"exitCode"`
	Error     string      `json:"error,omitempty"`
	StartedAt metav1.Time `json:"startedAt,omitempty"`
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}

type AppStatusTemplate struct {
	Stderr    string      `json:"stderr,omitempty"`
	ExitCode  int         `json:"exitCode"`
	Error     string      `json:"error,omitempty"`
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}

type AppStatusDeploy struct {
	Stdout    string      `json:"stdout,omitempty"`
	Stderr    string      `json:"stderr,omitempty"`
	Finished  bool        `json:"finished"`
	ExitCode  int         `json:"exitCode"`
	Error     string      `json:"error,omitempty"`
	StartedAt metav1.Time `json:"startedAt,omitempty"`
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}

type AppStatusInspect struct {
	Stdout    string      `json:"stdout,omitempty"`
	Stderr    string      `json:"stderr,omitempty"`
	ExitCode  int         `json:"exitCode"`
	Error     string      `json:"error,omitempty"`
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`
}
