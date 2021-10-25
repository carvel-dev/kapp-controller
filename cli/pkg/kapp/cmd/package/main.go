// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

// import (
// 	"os"

// 	"github.com/aunum/log"
// 	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

// 	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/component"
// )

// var descriptor = cliv1alpha1.PluginDescriptor{
// 	Name:        "package",
// 	Description: "Tanzu package management",
// 	Group:       cliv1alpha1.RunCmdGroup,
// }

// var logLevel int32
// var outputFormat string

// func main() {
// 	p, err := plugin.NewPlugin(&descriptor)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	p.Cmd.PersistentFlags().Int32VarP(&logLevel, "verbose", "", 0, "Number for the log level verbosity(0-9)")

// 	p.AddCommands(
// 		repositoryCmd,
// 		packageInstallCmd,
// 		packageAvailableCmd,
// 		packageInstalledCmd,
// 	)
// 	if err := p.Execute(); err != nil {
// 		os.Exit(1)
// 	}
// }

// // getOutputFormat gets the desired output format for package commands that need the ListTable format
// // for its output.
// func getOutputFormat() string {
// 	format := outputFormat
// 	if format != string(component.JSONOutputType) && format != string(component.YAMLOutputType) {
// 		// For table output, we want to force the list table format for this part
// 		format = string(component.ListTableOutputType)
// 	}
// 	return format
// }
