// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

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
		Short: "Trigger reconciliation for app",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "app", "a", "", "Set app name (required)")
	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   5 * time.Minute,
	})

	return cmd
}

func (o *KickOptions) Run() error {
	if len(o.Name) == 0 {
		return fmt.Errorf("Expected app name to be non empty")
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	app, err := client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("App '%s' does not exist in namespace '%s'", o.Name, o.NamespaceFlags.Name)
		}
		return err
	}

	if isOwnedByPackageInstall(app) {
		o.ui.BeginLinef("App '%s' is owned by '%s'\n", o.Name, fmt.Sprintf("%s/%s", app.OwnerReferences[0].Kind, app.OwnerReferences[0].Name))
	}

	o.ui.BeginLinef("Triggering reconciliation for app '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
	err = o.ui.AskForConfirmation()
	if err != nil {
		return err
	}

	err = o.triggerReconciliation(client)
	if err != nil {
		return err
	}

	if o.WaitFlags.Enabled {
		err = o.waitForAppReconciliation(client)
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

	o.statusUI.PrintMessagef("Triggering reconciliation for app '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)

	_, err = client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
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

	_, err = client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (o *KickOptions) waitForAppReconciliation(client kcclient.Interface) error {
	o.statusUI.PrintMessagef("Waiting for app reconciliation for '%s'", o.Name)
	appWatcher := NewAppTailer(o.NamespaceFlags.Name, o.Name, o.ui, client, AppTailerOpts{})

	err := appWatcher.TailAppStatus()
	if err != nil {
		return err
	}

	return nil
}
