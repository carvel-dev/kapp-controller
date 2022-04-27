// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	"github.com/spf13/cobra"
	kcpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "installed",
		Aliases: []string{"i"},
		Short:   "Manage installed packages",
	}
	return cmd
}

func packageInstallStatus(pkgi *kcpkgv1alpha1.PackageInstall) string {
	if pkgi.Spec.Canceled {
		return "Canceled"
	}
	if pkgi.Spec.Paused {
		return "Paused"
	}
	return pkgi.Status.FriendlyDescription
}
