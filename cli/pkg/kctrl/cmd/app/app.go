// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/spf13/cobra"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
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
