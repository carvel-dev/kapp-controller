// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"context"
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/k14s/kapp/pkg/kapp/cmd/core"
	"github.com/k14s/kapp/pkg/kapp/logger"
	"github.com/spf13/cobra"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type DeleteOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger
	pkgiName    string

	pollInterval time.Duration
	pollTimeout  time.Duration
	wait         bool

	nonInteractive bool

	NamespaceFlags cmdcore.NamespaceFlags
}

func NewDeleteOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *DeleteOptions {
	return &DeleteOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewDeleteCmd(o *DeleteOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"d"},
		Short:   "Uninstall installed Package",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVar(&o.pkgiName, "name", "", "Name of PackageInstall")
	cmd.Flags().DurationVar(&o.pollInterval, "poll-interval", 1*time.Second, "Time interval between consecutive polls while reconciling")
	cmd.Flags().DurationVar(&o.pollTimeout, "poll-timeout", 1*time.Minute, "Timeout for the reconciliation process")
	cmd.Flags().BoolVar(&o.wait, "wait", true, "Wait for reconcilation, default true")
	cmd.Flags().BoolVarP(&o.nonInteractive, "yes", "y", false, "Do not ask for confirmation. default false")
	return cmd
}

func (o *DeleteOptions) Run() error {
	o.ui.PrintLinef("Deleting package install '%s' from namespace '%s'", o.pkgiName, o.NamespaceFlags.Name)

	if !o.nonInteractive {
		err := o.ui.AskForConfirmation()
		if err != nil {
			return err
		}
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return nil
	}

	err = client.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Delete(
		context.Background(), o.pkgiName, metav1.DeleteOptions{},
	)
	if err != nil {
		return err
	}

	if !o.wait {
		return nil
	}

	o.ui.PrintLinef("Waiting for deletion of PackageInstall '%s' from namespace '%s'", o.pkgiName, o.NamespaceFlags.Name)
	err = o.waitForResourceDelete()
	if err != nil {
		return err
	}

	return nil

	//TODO: Handle "created resources" (Secret, ClusterRoleBinding, ClusterRole, ServiceAccount)
}

func (o *DeleteOptions) waitForResourceDelete() error {
	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	err = wait.Poll(o.pollInterval, o.pollTimeout, func() (bool, error) {
		resource, err := client.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Get(
			context.Background(), o.pkgiName, metav1.GetOptions{},
		)
		if err != nil {
			if errors.IsNotFound(err) {
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
			o.ui.PrintLinef("'PackageInstall' resource deletion status: %s", cond.Type)
			if cond.Type == kcv1alpha1.DeleteFailed && cond.Status == corev1.ConditionTrue {
				return false, fmt.Errorf("resource deletion failed: %s. %s", status.UsefulErrorMessage, status.FriendlyDescription)
			}
		}
		return false, nil
	})

	if err != nil {
		return err
	}

	return nil
}
