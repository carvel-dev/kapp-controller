// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

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
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type KickOptions struct {
	ui          ui.UI
	statusUI    cmdcore.StatusLoggingUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	WaitFlags cmdcore.WaitFlags

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string
}

func NewKickOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *KickOptions {
	return &KickOptions{ui: ui, statusUI: cmdcore.NewStatusLoggingUI(ui), depsFactory: depsFactory, logger: logger}
}

func NewKickCmd(o *KickOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kick",
		Short: "Trigger reconciliation for repository",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
		Annotations: map[string]string{cmdapp.TTYByDefaultKey: "",
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value},
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "repository", "r", "", "Set repository name (required)")
	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   5 * time.Minute,
	})

	return cmd
}

func (o *KickOptions) Run() error {
	if len(o.Name) == 0 {
		return fmt.Errorf("Expected repository name to be non empty")
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	_, err = client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("Repository '%s' does not exist in namespace '%s'", o.Name, o.NamespaceFlags.Name)
		}
		return err
	}

	o.ui.BeginLinef("Triggering reconciliation for package repository '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
	err = o.ui.AskForConfirmation()
	if err != nil {
		return err
	}

	err = o.triggerReconciliation(client)
	if err != nil {
		return err
	}

	if o.WaitFlags.Enabled {
		err = o.waitForReconciliation(client)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *KickOptions) triggerReconciliation(client kcclient.Interface) error {
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

	o.statusUI.PrintMessagef("Triggering reconciliation for package repository '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)

	_, err = client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	unpausePatch := []map[string]interface{}{
		{
			"op":   "remove",
			"path": "/spec/paused",
		},
	}

	patchJSON, err = json.Marshal(unpausePatch)
	if err != nil {
		return err
	}

	_, err = client.PackagingV1alpha1().PackageRepositories(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (o *KickOptions) waitForReconciliation(client kcclient.Interface) error {
	o.statusUI.PrintMessagef("Waiting for package repository reconciliation for '%s'", o.Name)
	repoWatcher := NewRepoTailer(o.NamespaceFlags.Name, o.Name, o.ui, client)

	err := repoWatcher.TailRepoStatus()
	if err != nil {
		return err
	}

	return nil
}
