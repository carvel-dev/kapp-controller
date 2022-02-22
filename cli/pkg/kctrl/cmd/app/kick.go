// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type KickOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string
}

func NewKickOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *KickOptions {
	return &KickOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewKickCmd(o *KickOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kick",
		Short: "Trigger reconciliation for App CR",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "app", "a", "", "Set App CR name (required)")

	return cmd
}

func (o *KickOptions) Run() error {
	if len(o.Name) == 0 {
		return fmt.Errorf("Expected App CR name to be non empty")
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	app, err := client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("App CR '%s' does not exist in namespace '%s'", o.Name, o.NamespaceFlags.Name)
		}
		return err
	}

	if isOwnedByPackageInstall(app) {
		return fmt.Errorf("App CR '%s' in namespace '%s' is owned by a PackageInstall.\n(Hint: Try using `kctrl package installed kick` to reconcile PackageInstall)",
			o.Name, o.NamespaceFlags.Name)
	}

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

	o.ui.PrintLinef("Triggering reconciliation for App CR '%s' in namespace '%s'...", o.Name, o.NamespaceFlags.Name)

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
