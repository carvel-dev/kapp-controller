// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdapp "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type DeleteOptions struct {
	ui          ui.UI
	statusUI    cmdcore.StatusLoggingUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	WaitFlags cmdcore.WaitFlags

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewDeleteOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *DeleteOptions {
	return &DeleteOptions{ui: ui, statusUI: cmdcore.NewStatusLoggingUI(ui), depsFactory: depsFactory, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewDeleteCmd(o *DeleteOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a package repository",
		RunE:  func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Delete a package repository",
				[]string{"package", "repository", "delete", "-r", "sample-repo"}},
		}.Description("-r", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations: map[string]string{cmdapp.TTYByDefaultKey: "",
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value},
	}

	o.NamespaceFlags.SetWithPackageCommandTreeOpts(cmd, flagsFactory, o.pkgCmdTreeOpts)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "repository", "r", "", "Set package repository name (required)")
	} else {
		cmd.Use = "delete REPOSITORY_NAME"
		cmd.Args = cobra.ExactArgs(1)
	}

	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   5 * time.Minute,
	})

	return cmd
}

func (o *DeleteOptions) Run(args []string) error {
	if o.pkgCmdTreeOpts.PositionalArgs {
		if len(args) > 0 {
			o.Name = args[0]
		}
	}

	if len(o.Name) == 0 {
		return fmt.Errorf("Expected package repository name to be non-empty")
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	o.ui.PrintLinef("Deleting package repository '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)

	err = o.ui.AskForConfirmation()
	if err != nil {
		return err
	}

	err = client.PackagingV1alpha1().PackageRepositories(
		o.NamespaceFlags.Name).Delete(context.Background(), o.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	if o.WaitFlags.Enabled {
		return o.waitForDeletion(client)
	}
	return nil
}

func (o *DeleteOptions) waitForDeletion(client versioned.Interface) error {
	o.statusUI.PrintMessagef("Waiting for package repository reconciliation for '%s'", o.Name)

	msgsUI := cmdcore.NewDedupingMessagesUI(cmdcore.NewPlainMessagesUI(o.ui))
	description := getPackageRepositoryDescription(o.Name, o.NamespaceFlags.Name)

	repoStatusTailErrored := false
	tailRepoStatusOutput := func(tailErrored *bool) {
		repoWatcher := NewRepoTailer(o.NamespaceFlags.Name, o.Name, o.ui, client, RepoTailerOpts{})

		err := repoWatcher.TailRepoStatus()
		if err != nil {
			o.statusUI.PrintMessagef("Error tailing or reconciling Package Repository: %s", err.Error())
			*tailErrored = true
		}
	}

	go tailRepoStatusOutput(&repoStatusTailErrored)

	if err := wait.Poll(o.WaitFlags.CheckInterval, o.WaitFlags.Timeout, func() (bool, error) {
		resource, err := client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
		if err != nil {
			if !(errors.IsNotFound(err)) {
				return true, nil
			}
		}
		if err != nil {
			if errors.IsNotFound(err) {
				msgsUI.NotifySection("%s: DeletionSucceeded", description)
				return true, nil
			}
			return false, err
		}

		if resource.Generation != resource.Status.ObservedGeneration {
			// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
			return false, nil
		}
		status := resource.Status.GenericStatus

		for _, cond := range status.Conditions {
			if repoStatusTailErrored {
				msgsUI.NotifySection("%s: %s", description, cond.Type)
			}

			if cond.Type == kcv1alpha1.DeleteFailed && cond.Status == corev1.ConditionTrue {
				return false, fmt.Errorf("%s: Deleting: %s. %s", description, status.UsefulErrorMessage, status.FriendlyDescription)
			}
		}
		return false, nil
	}); err != nil {
		return fmt.Errorf("%s: Deleting: %s", description, err)
	}
	return nil
}
