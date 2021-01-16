package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type GenericStatus struct {
	ObservedGeneration  int64          `json:"observedGeneration"`
	Conditions          []AppCondition `json:"conditions"`
	FriendlyDescription string         `json:"friendlyDescription"`
}

type AppConditionType string

const (
	Reconciling        AppConditionType = "Reconciling"
	ReconcileFailed    AppConditionType = "ReconcileFailed"
	ReconcileSucceeded AppConditionType = "ReconcileSucceeded"

	Deleting     AppConditionType = "Deleting"
	DeleteFailed AppConditionType = "DeleteFailed"
)

// TODO rename to Condition
type AppCondition struct {
	Type   AppConditionType       `json:"type"`
	Status corev1.ConditionStatus `json:"status"`
	// Unique, this should be a short, machine understandable string that gives the reason
	// for condition's last transition. If it reports "ResizeStarted" that means the underlying
	// persistent volume is being resized.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}
