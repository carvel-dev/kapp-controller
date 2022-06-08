// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/spf13/cobra"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "app",
		Aliases: []string{"a"},
		Short:   "App",
	}
	return cmd
}

func isOwnedByPackageInstall(app *kcv1alpha1.App) bool {
	for _, reference := range app.OwnerReferences {
		if reference.APIVersion == kcpkgv1alpha1.SchemeGroupVersion.Identifier() {
			return true
		}
	}
	return false
}

func appStatusString(app *kcv1alpha1.App) string {
	if len(app.Status.Conditions) < 1 {
		return ""
	}
	if app.Spec.Canceled {
		return "Canceled"
	}
	if app.Spec.Paused {
		return "Paused"
	}
	for _, condition := range app.Status.Conditions {
		switch condition.Type {
		case kcv1alpha1.ReconcileFailed:
			return "Reconcile failed"
		case kcv1alpha1.ReconcileSucceeded:
			return "Reconcile succeeded"
		case kcv1alpha1.DeleteFailed:
			return "Deletion failed"
		case kcv1alpha1.Reconciling:
			return "Reconciling"
		case kcv1alpha1.Deleting:
			return "Deleting"
		}
	}
	return app.Status.FriendlyDescription
}

func isFailing(conditions []kcv1alpha1.AppCondition) bool {
	for _, condition := range conditions {
		if condition.Type == kcv1alpha1.ReconcileFailed && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}
