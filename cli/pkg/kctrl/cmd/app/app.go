// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"

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

func getAppDescription(name string, namespace string) string {
	description := fmt.Sprintf("app/%s (kappctrl.k14s.io/v1alpha1)  namespace: %s", name, namespace)
	return description
}
