// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

// import (
// 	"fmt"

// 	"github.com/spf13/cobra"

// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackageclient"
// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
// )

// var repositoryUpdateCmd = &cobra.Command{
// 	Use:   "update REPOSITORY_NAME --url REPOSITORY_URL",
// 	Short: "Update a package repository",
// 	Args:  cobra.ExactArgs(1),
// 	Example: `
//     # Update repository in default namespace
//     tanzu package repository update repo --url projects-stg.registry.vmware.com/tkg/standard-repo:v1.0.1 --namespace test-ns`,
// 	RunE:         repositoryUpdate,
// 	SilenceUsage: true,
// }

// func init() {
// 	repositoryUpdateCmd.Flags().StringVarP(&repoOp.RepositoryURL, "url", "", "", "OCI registry url for package repository bundle")
// 	repositoryUpdateCmd.Flags().BoolVarP(&repoOp.CreateRepository, "create", "", false, "Creates the package repository if it does not exist, optional")
// 	repositoryUpdateCmd.Flags().BoolVarP(&repoOp.CreateNamespace, "create-namespace", "", false, "Create namespace if the target namespace does not exist, optional")
// 	repositoryUpdateCmd.Flags().BoolVarP(&repoOp.Wait, "wait", "", true, "Wait for the package repository reconciliation to complete, optional. To disable wait, specify --wait=false")
// 	repositoryUpdateCmd.Flags().DurationVarP(&repoOp.PollInterval, "poll-interval", "", tkgpackagedatamodel.DefaultPollInterval, "Time interval between subsequent polls of package repository reconciliation status, optional")
// 	repositoryUpdateCmd.Flags().DurationVarP(&repoOp.PollTimeout, "poll-timeout", "", tkgpackagedatamodel.DefaultPollTimeout, "Timeout value for polls of package repository reconciliation status, optional")
// 	repositoryUpdateCmd.MarkFlagRequired("url") //nolint
// 	repositoryCmd.AddCommand(repositoryUpdateCmd)
// }

// func repositoryUpdate(cmd *cobra.Command, args []string) error {
// 	repoOp.RepositoryName = args[0]

// 	pkgClient, err := tkgpackageclient.NewTKGPackageClient(repoOp.KubeConfig)
// 	if err != nil {
// 		return err
// 	}

// 	pp := &tkgpackagedatamodel.PackageProgress{
// 		ProgressMsg: make(chan string, 10),
// 		Err:         make(chan error),
// 		Done:        make(chan struct{}),
// 	}
// 	go pkgClient.UpdateRepository(repoOp, pp)

// 	initialMsg := fmt.Sprintf("Updating package repository '%s'", repoOp.RepositoryName)
// 	if err := DisplayProgress(initialMsg, pp); err != nil {
// 		if err.Error() == tkgpackagedatamodel.ErrRepoNotExists {
// 			log.Warningf("\npackage repository '%s' does not exist in namespace '%s'. Consider using the --create flag to add the package repository", repoOp.RepositoryName, repoOp.Namespace)
// 			return nil
// 		}
// 		return err
// 	}

// 	log.Infof("\n Updated package repository '%s' in namespace '%s'", repoOp.RepositoryName, repoOp.Namespace)
// 	return nil
// }
