// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

// import (
// 	"github.com/spf13/cobra"

// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
// )

// var packageInstalledCreateCmd = &cobra.Command{
// 	Use:          "create INSTALLED_PACKAGE_NAME --package-name PACKAGE_NAME --version VERSION",
// 	Short:        "Install a package",
// 	Args:         cobra.ExactArgs(1),
// 	RunE:         packageInstall,
// 	SilenceUsage: true,
// }

// func init() {
// 	packageInstalledCreateCmd.Flags().StringVarP(&packageInstallOp.PackageName, "package-name", "p", "", "Name of the package to be installed")
// 	packageInstalledCreateCmd.Flags().StringVarP(&packageInstallOp.Version, "version", "v", "", "Version of the package to be installed")
// 	packageInstalledCreateCmd.Flags().BoolVarP(&packageInstallOp.CreateNamespace, "create-namespace", "", false, "Create namespace if the target namespace does not exist, optional")
// 	packageInstalledCreateCmd.Flags().StringVarP(&packageInstallOp.Namespace, "namespace", "n", "default", "Target namespace to install the package, optional")
// 	packageInstalledCreateCmd.Flags().StringVarP(&packageInstallOp.ServiceAccountName, "service-account-name", "", "", "Name of an existing service account used to install underlying package contents, optional")
// 	packageInstalledCreateCmd.Flags().StringVarP(&packageInstallOp.ValuesFile, "values-file", "f", "", "The path to the configuration values file, optional")
// 	packageInstalledCreateCmd.Flags().BoolVarP(&packageInstallOp.Wait, "wait", "", true, "Wait for the package reconciliation to complete, optional")
// 	packageInstalledCreateCmd.Flags().DurationVarP(&packageInstallOp.PollInterval, "poll-interval", "", tkgpackagedatamodel.DefaultPollInterval, "Time interval between subsequent polls of package reconciliation status, optional")
// 	packageInstalledCreateCmd.Flags().DurationVarP(&packageInstallOp.PollTimeout, "poll-timeout", "", tkgpackagedatamodel.DefaultPollTimeout, "Timeout value for polls of package reconciliation status, optional")
// 	packageInstalledCreateCmd.PersistentFlags().StringVarP(&packageInstallOp.KubeConfig, "kubeconfig", "", "", "The path to the kubeconfig file, optional")
// 	packageInstalledCreateCmd.MarkFlagRequired("package-name") //nolint
// 	packageInstalledCreateCmd.MarkFlagRequired("version")      //nolint
// 	packageInstalledCmd.AddCommand(packageInstalledCreateCmd)
// }
