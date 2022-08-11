// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reconciler

import (
	"fmt"
	"strings"

	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Status struct {
	S          kcv1alpha1.GenericStatus
	UpdateFunc func(kcv1alpha1.GenericStatus)
}

func (s *Status) Result() kcv1alpha1.GenericStatus { return s.S }

func (s *Status) IsReconciling() bool {
	for _, cond := range s.S.Conditions {
		if cond.Type == kcv1alpha1.Reconciling {
			return true
		}
	}
	return false
}

func (s *Status) IsReconcileSucceeded() bool {
	for _, cond := range s.S.Conditions {
		if cond.Type == kcv1alpha1.ReconcileSucceeded {
			return true
		}
	}
	return false
}

func (s *Status) IsReconcileFailed() bool {
	for _, cond := range s.S.Conditions {
		if cond.Type == kcv1alpha1.ReconcileFailed {
			return true
		}
	}
	return false
}

func (s *Status) IsDeleting() bool {
	for _, cond := range s.S.Conditions {
		if cond.Type == kcv1alpha1.Deleting {
			return true
		}
	}
	return false
}

func (s *Status) IsDeleteFailed() bool {
	for _, cond := range s.S.Conditions {
		if cond.Type == kcv1alpha1.DeleteFailed {
			return true
		}
	}
	return false
}

func (s *Status) SetReconciling(meta metav1.ObjectMeta) {
	s.markObservedLatest(meta)
	s.removeAllConditions()

	s.S.Conditions = append(s.S.Conditions, kcv1alpha1.Condition{
		Type:   kcv1alpha1.Reconciling,
		Status: corev1.ConditionTrue,
	})

	s.S.FriendlyDescription = "Reconciling"
	s.S.UsefulErrorMessage = ""

	s.UpdateFunc(s.S)
}

func (s *Status) SetDeleting(meta metav1.ObjectMeta) {
	s.markObservedLatest(meta)
	s.removeAllConditions()

	s.S.Conditions = append(s.S.Conditions, kcv1alpha1.Condition{
		Type:   kcv1alpha1.Deleting,
		Status: corev1.ConditionTrue,
	})

	s.S.FriendlyDescription = "Deleting"

	s.UpdateFunc(s.S)
}

func (s *Status) SetReconcileCompleted(err error) {
	s.removeAllConditions()

	if err != nil {
		s.S.Conditions = append(s.S.Conditions, kcv1alpha1.Condition{
			Type:    kcv1alpha1.ReconcileFailed,
			Status:  corev1.ConditionTrue,
			Message: err.Error(),
		})
		s.S.FriendlyDescription = s.friendlyErrMsg(fmt.Sprintf("Reconcile failed: %s", err.Error()))
	} else {
		s.S.Conditions = append(s.S.Conditions, kcv1alpha1.Condition{
			Type:    kcv1alpha1.ReconcileSucceeded,
			Status:  corev1.ConditionTrue,
			Message: "",
		})
		s.S.FriendlyDescription = "Reconcile succeeded"
		s.S.UsefulErrorMessage = ""
	}

	s.UpdateFunc(s.S)
}

func (s *Status) SetDeleteCompleted(err error) {
	s.removeAllConditions()

	if err != nil {
		s.S.Conditions = append(s.S.Conditions, kcv1alpha1.Condition{
			Type:    kcv1alpha1.DeleteFailed,
			Status:  corev1.ConditionTrue,
			Message: err.Error(),
		})
		s.S.FriendlyDescription = s.friendlyErrMsg(fmt.Sprintf("Delete failed: %s", err.Error()))
	} else {
		// assume resource will be deleted, hence nothing to update
	}

	s.UpdateFunc(s.S)
}

func (s *Status) friendlyErrMsg(errMsg string) string {
	errMsgPieces := strings.Split(errMsg, "\n")
	if len(errMsgPieces[0]) > 80 {
		errMsgPieces[0] = errMsgPieces[0][:80] + "..."
	} else if len(errMsgPieces) > 1 {
		errMsgPieces[0] += "..."
	}
	return errMsgPieces[0]
}

func (s *Status) markObservedLatest(meta metav1.ObjectMeta) {
	s.S.ObservedGeneration = meta.Generation
}

func (s *Status) removeAllConditions() {
	s.S.Conditions = nil
}

func (s *Status) SetUsefulErrorMessage(errMsg string) {
	s.S.UsefulErrorMessage = errMsg
}
