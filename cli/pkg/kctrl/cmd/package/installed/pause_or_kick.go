// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type PauseOrKickOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	WaitFlags cmdcore.WaitFlags

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewPauseOrKickOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *PauseOrKickOptions {
	return &PauseOrKickOptions{ui: ui, depsFactory: depsFactory, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewPauseCmd(o *PauseOrKickOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause reconciliation of package install",
		Args:  cobra.ExactArgs(1),
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Pause() },
		Example: cmdcore.Examples{
			cmdcore.Example{"Pause reconciliation of package install",
				[]string{"package", "installed", "pause", "-i", "cert-man"},
			},
		}.Description("-i", o.pkgCmdTreeOpts),
		SilenceUsage: true,
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "package-install", "i", "", "Set installed package name (required)")
	} else {
		cmd.Use = "pause INSTALLED_PACKAGE_NAME"
	}

	return cmd
}

func NewKickCmd(o *PauseOrKickOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kick",
		Short: "Trigger reconciliation of package install",
		Args:  cobra.ExactArgs(1),
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Kick() },
		Example: cmdcore.Examples{
			cmdcore.Example{"Trigger reconciliation of package install",
				[]string{"package", "installed", "kick", "-i", "cert-man"},
			},
		}.Description("-i", o.pkgCmdTreeOpts),
		SilenceUsage: true,
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "package-install", "i", "", "Set installed package name (required)")
	} else {
		cmd.Use = "kick INSTALLED_PACKAGE_NAME"
	}

	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   5 * time.Minute,
	})

	return cmd
}

func (o *PauseOrKickOptions) Pause() error {
	if len(o.Name) == 0 {
		return fmt.Errorf("Expected package install name to be non empty")
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	_, err = client.PackagingV1alpha1().PackageInstalls(
		o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	return o.pause(client)
}

func (o *PauseOrKickOptions) Kick() error {
	if len(o.Name) == 0 {
		return fmt.Errorf("Expected package install name to be non empty")
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	_, err = client.PackagingV1alpha1().PackageInstalls(
		o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	err = o.pause(client)
	if err != nil {
		return err
	}

	return o.unpause(client)
}

func (o *PauseOrKickOptions) pause(client kcclient.Interface) error {
	pausePatch := []map[string]interface{}{
		{
			"op":    "add",
			"path":  "/spec/paused",
			"value": true,
		},
	}

	patchJSON, err := json.Marshal(pausePatch)
	if err != nil {
		return err
	}

	o.ui.PrintLinef("Pausing reconciliation for package install '%s' in namespace '%s'...", o.Name, o.NamespaceFlags.Name)

	_, err = client.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (o *PauseOrKickOptions) unpause(client kcclient.Interface) error {
	unpausePatch := []map[string]interface{}{
		{
			"op":   "remove",
			"path": "/spec/paused",
		},
	}

	patchJSON, err := json.Marshal(unpausePatch)
	if err != nil {
		return err
	}

	o.ui.PrintLinef("Triggering reconciliation for package install '%s' in namespace '%s'...", o.Name, o.NamespaceFlags.Name)

	_, err = client.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}
