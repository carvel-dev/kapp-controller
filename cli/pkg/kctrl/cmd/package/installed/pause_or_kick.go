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
	cmdapp "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

type PauseOrKickOptions struct {
	ui          ui.UI
	statusUI    cmdcore.StatusLoggingUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	WaitFlags cmdcore.WaitFlags

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewPauseOrKickOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *PauseOrKickOptions {
	return &PauseOrKickOptions{ui: ui, statusUI: cmdcore.NewStatusLoggingUI(ui), depsFactory: depsFactory, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
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

	o.ui.PrintLinef("Pausing reconciliation for package install '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
	err = o.ui.AskForConfirmation()
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

	o.ui.PrintLinef("Triggering reconciliation for package install '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
	err = o.ui.AskForConfirmation()
	if err != nil {
		return err
	}

	err = o.pause(client)
	if err != nil {
		return err
	}

	err = o.unpause(client)
	if err != nil {
		return err
	}

	if o.WaitFlags.Enabled {
		return o.waitForPackageInstallReconciliation(client)
	}

	return nil
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

	_, err = client.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

// waitForPackageInstallReconciliation waits until the package get installed successfully or a failure happen
// TODO Move reconciliation to a common place for create-or-update and pause-or-kick
func (o *PauseOrKickOptions) waitForPackageInstallReconciliation(client kcclient.Interface) error {
	o.statusUI.PrintMessagef("Waiting for PackageInstall reconciliation for '%s'", o.Name)
	msgsUI := cmdcore.NewDedupingMessagesUI(cmdcore.NewPlainMessagesUI(o.ui))
	description := getPackageInstallDescription(o.Name, o.NamespaceFlags.Name)

	appStatusTailErrored := false
	tailAppStatusOutput := func(tailErrored *bool) {
		appWatcher := cmdapp.NewAppTailer(o.NamespaceFlags.Name, o.Name, o.ui, client, cmdapp.AppTailerOpts{
			IgnoreNotExists: true,
		})

		err := appWatcher.TailAppStatus()
		if err != nil {
			o.statusUI.PrintMessagef("Error tailing or reconciling app: %s", err.Error())
			*tailErrored = true
		}
	}
	go tailAppStatusOutput(&appStatusTailErrored)

	if err := wait.Poll(o.WaitFlags.CheckInterval, o.WaitFlags.Timeout, func() (done bool, err error) {
		resource, err := client.PackagingV1alpha1().PackageInstalls(o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if resource.Generation != resource.Status.ObservedGeneration {
			// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
			return false, nil
		}
		status := resource.Status.GenericStatus

		for _, condition := range status.Conditions {
			if appStatusTailErrored {
				msgsUI.NotifySection("%s: %s", description, condition.Type)
			}

			switch {
			case condition.Type == kcv1alpha1.ReconcileSucceeded && condition.Status == corev1.ConditionTrue:
				return true, nil
			case condition.Type == kcv1alpha1.ReconcileFailed && condition.Status == corev1.ConditionTrue:
				return false, fmt.Errorf("%s. %s", status.UsefulErrorMessage, status.FriendlyDescription)
			}
		}
		return false, nil
	}); err != nil {
		return fmt.Errorf("%s: Reconciling: %s", description, err)
	}

	return nil
}
