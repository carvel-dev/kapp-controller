// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeleteOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	WaitFlags cmdcore.WaitFlags

	positionalNameArg bool
}

func NewDeleteOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, positionalNameArg bool) *DeleteOptions {
	return &DeleteOptions{ui: ui, depsFactory: depsFactory, logger: logger, positionalNameArg: positionalNameArg}
}

func NewDeleteCmd(o *DeleteOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a package repository",
		RunE:  func(_ *cobra.Command, args []string) error { return o.Run(args) },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)

	if !o.positionalNameArg {
		cmd.Flags().StringVarP(&o.Name, "repository", "r", "", "Set package repository name")
	}

	o.WaitFlags.Set(cmd, flagsFactory)

	return cmd
}

func (o *DeleteOptions) Run(args []string) error {
	if o.positionalNameArg {
		o.Name = args[0]
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
	o.ui.PrintLinef("Waiting for deletion to be completed...")

	t1 := time.Now()

	for {
		pkgr, err := client.PackagingV1alpha1().PackageRepositories(
			o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				o.ui.PrintLinef("Repository deleted successfully")
				return nil
			}
			return err
		}

		// Should wait for generation to be observed before checking
		// the reconciliation status so that we know we are checking the new spec
		if pkgr.Generation == pkgr.Status.ObservedGeneration {
			for _, condition := range pkgr.Status.Conditions {
				o.ui.BeginLinef("'PackageRepository' resource deletion status: %s\n", condition.Type)
				if condition.Type == v1alpha1.DeleteFailed && condition.Status == corev1.ConditionTrue {
					return fmt.Errorf("Deletion failed: %s", pkgr.Status.UsefulErrorMessage)
				}
			}
		}

		if time.Now().Sub(t1) > o.WaitFlags.Timeout {
			return fmt.Errorf("Timed out waiting for package repository to be deleted")
		}

		time.Sleep(o.WaitFlags.CheckInterval)
	}
}
