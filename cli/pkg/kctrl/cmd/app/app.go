// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"

	"github.com/cppforlife/color"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "app",
		Aliases: []string{"a"},
		Short:   "App",
		Annotations: map[string]string{
			cmdcore.AppHelpGroup.Key: cmdcore.AppHelpGroup.Value,
		},
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

// Returns app status string and a bool indicating if it is a failure
func appStatus(app *kcv1alpha1.App) (string, bool) {
	if len(app.Status.Conditions) == 0 {
		return "", false
	}
	if app.Spec.Canceled {
		return "Canceled", true
	}
	if app.Spec.Paused {
		return "Paused", true
	}
	for _, condition := range app.Status.Conditions {
		switch condition.Type {
		case kcv1alpha1.ReconcileFailed:
			return "Reconcile failed", true
		case kcv1alpha1.ReconcileSucceeded:
			return "Reconcile succeeded", false
		case kcv1alpha1.DeleteFailed:
			return "Deletion failed", true
		case kcv1alpha1.Reconciling:
			return "Reconciling", false
		case kcv1alpha1.Deleting:
			return "Deleting", false
		}
	}
	return app.Status.FriendlyDescription, false
}

func HasReconciled(status kcv1alpha1.AppStatus) bool {
	for _, condition := range status.Conditions {
		if condition.Type == kcv1alpha1.ReconcileSucceeded && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func HasFailed(status kcv1alpha1.AppStatus) (bool, string) {
	for _, condition := range status.Conditions {
		if condition.Type == kcv1alpha1.ReconcileFailed && condition.Status == corev1.ConditionTrue {
			return true, color.RedString(fmt.Sprintf("%s: %s", kcv1alpha1.ReconcileFailed, status.UsefulErrorMessage))
		}
		if condition.Type == kcv1alpha1.DeleteFailed && condition.Status == corev1.ConditionTrue {
			return true, color.RedString(fmt.Sprintf("%s: %s", kcv1alpha1.DeleteFailed, status.UsefulErrorMessage))
		}
	}
	return false, ""
}

func IsDeleting(status kcv1alpha1.AppStatus) bool {
	for _, condition := range status.Conditions {
		if condition.Type == kcv1alpha1.Deleting && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}
