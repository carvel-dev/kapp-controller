// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

// import (
// 	"github.com/spf13/cobra"

// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
// )

// var packageAvailableOp = tkgpackagedatamodel.NewPackageAvailableOptions()

// var packageAvailableCmd = &cobra.Command{
// 	Use:       "available",
// 	ValidArgs: []string{"list", "get"},
// 	Short:     "Manage available packages",
// 	Args:      cobra.RangeArgs(1, 2),
// }

// func init() {
// 	packageAvailableCmd.PersistentFlags().StringVarP(&packageAvailableOp.KubeConfig, "kubeconfig", "", "", "The path to the kubeconfig file, optional")
// 	packageAvailableCmd.PersistentFlags().StringVarP(&packageAvailableOp.Namespace, "namespace", "n", "default", "Namespace of packages, optional")
// 	packageAvailableCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml|json|table), optional")
// }
