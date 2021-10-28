// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"context"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/k14s/kapp/pkg/kapp/cmd/core"

	"github.com/k14s/kapp/pkg/kapp/logger"
	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeleteOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags  cmdcore.NamespaceFlags
	RepositoryName  string
	Wait            bool
	NonIntereactive bool
	PollInterval    time.Duration
	PollTimeout     time.Duration
}

func NewDeleteOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *DeleteOptions {
	return &DeleteOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewDeleteCmd(o *DeleteOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"l", "ls"},
		Short:   "Delete a package repository",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.RepositoryName, "repository", "R", "", "Delete a package repository")
	cmd.Flags().BoolVarP(&o.Wait, "wait", "", true, "Wait for the package repository reconciliation to complete, optional. To disable wait, specify --wait=false")
	cmd.Flags().BoolVarP(&o.NonIntereactive, "yes", "y", false, "Delete package repository without asking for confirmation, optional")
	cmd.Flags().DurationVarP(&o.PollInterval, "poll-interval", "", 1*time.Second, "Time interval between subsequent polls of package repository reconciliation status, optional")
	cmd.Flags().DurationVarP(&o.PollTimeout, "poll-timeout", "", 5*time.Minute, "Timeout value for polls of package repository reconciliation status, optional")
	return cmd
}

func (o *DeleteOptions) Run() error {
	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	o.ui.PrintLinef("Deleting package repository '%s' in namespace '%s'", o.RepositoryName, o.NamespaceFlags.Name)

	//TODO: Use ui_flags
	if !o.NonIntereactive {
		if err := o.ui.AskForConfirmation(); err != nil {
			return err
		}
	}

	err = client.PackagingV1alpha1().PackageRepositories(
		o.NamespaceFlags.Name).Delete(context.Background(), o.RepositoryName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	//TODO: Use better approach for waiting
	if o.Wait {
		start := time.Now()
		o.ui.PrintLinef("Waiting for deletion to be completed")
		for {
			_, err := client.PackagingV1alpha1().PackageRepositories(
				o.NamespaceFlags.Name).Get(context.Background(), o.RepositoryName, metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					break
				}
				return err
			}

			if time.Now().Before(start.Add(o.PollTimeout)) {
				break
			}
			time.Sleep(o.PollInterval)
		}
	}

	return nil
}
