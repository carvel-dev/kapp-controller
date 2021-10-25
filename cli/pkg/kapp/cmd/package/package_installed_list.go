// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

// import (
// 	"github.com/spf13/cobra"

// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/component"
// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/kappclient"
// )

// var packageInstalledListCmd = &cobra.Command{
// 	Use:   "list",
// 	Short: "List installed packages",
// 	Args:  cobra.NoArgs,
// 	Example: `
//     # List installed packages across all namespaces
//     tanzu package installed list -A

//     # List installed packages from specified namespace
//     tanzu package installed list --namespace test-ns`,
// 	RunE:         packageInstalledList,
// 	SilenceUsage: true,
// }

// func init() {
// 	packageInstalledListCmd.Flags().BoolVarP(&packageInstalledOp.AllNamespaces, "all-namespaces", "A", false, "If present, list packages across all namespaces, optional")
// 	packageInstalledListCmd.Flags().StringVarP(&packageInstalledOp.Namespace, "namespace", "n", "default", "Namespace for installed package CR, optional")
// 	packageInstalledListCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml|json|table), optional")
// 	packageInstalledCmd.AddCommand(packageInstalledListCmd)
// }

// func packageInstalledList(cmd *cobra.Command, args []string) error {
// 	kc, err := kappclient.NewKappClient(packageInstalledOp.KubeConfig)
// 	if err != nil {
// 		return err
// 	}
// 	if packageInstalledOp.AllNamespaces {
// 		packageInstalledOp.Namespace = ""
// 	}
// 	t, err := component.NewOutputWriterWithSpinner(cmd.OutOrStdout(), outputFormat,
// 		"Retrieving installed packages...", true)
// 	if err != nil {
// 		return err
// 	}

// 	pkgInstalledList, err := kc.ListPackageInstalls(packageInstalledOp.Namespace)
// 	if err != nil {
// 		return err
// 	}

// 	if packageInstalledOp.AllNamespaces {
// 		t.SetKeys("NAME", "PACKAGE-NAME", "PACKAGE-VERSION", "STATUS", "NAMESPACE")
// 	} else {
// 		t.SetKeys("NAME", "PACKAGE-NAME", "PACKAGE-VERSION", "STATUS")
// 	}
// 	for i := range pkgInstalledList.Items {
// 		pkg := pkgInstalledList.Items[i]
// 		if packageInstalledOp.AllNamespaces {
// 			t.AddRow(pkg.Name, pkg.Spec.PackageRef.RefName, pkg.Status.Version,
// 				pkg.Status.FriendlyDescription, pkg.Namespace)
// 		} else {
// 			t.AddRow(pkg.Name, pkg.Spec.PackageRef.RefName, pkg.Status.Version,
// 				pkg.Status.FriendlyDescription)
// 		}
// 	}
// 	t.RenderWithSpinner()
// 	return nil
// }
