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
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type PauseOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string
}

func NewPauseOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *PauseOptions {
	return &PauseOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewPauseCmd(o *PauseOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pause",
		Aliases: []string{"p"},
		Short:   "Pause reconciliation for App CR",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "app", "a", "", "Set App CR name (required)")

	return cmd
}

func (o *PauseOptions) Run() error {
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
		return fmt.Errorf("App CR '%s' in namespace '%s' is owned by a PackageInstall.\n(Hint: Try using `kctrl package installed pause` to pause PackageInstall)",
			o.Name, o.NamespaceFlags.Name)
	}

	if app.Spec.Paused {
		o.ui.PrintLinef("App CR '%s' in namespace '%s' is already paused", o.Name, o.NamespaceFlags.Name)
		return nil
	}

	err = o.pauseApp(client)
	if err != nil {
		return err
	}

	return nil
}

func (o *PauseOptions) pauseApp(client kcclient.Interface) error {
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

	o.ui.PrintLinef("Pausing reconciliation for App CR '%s' in namespace '%s'...", o.Name, o.NamespaceFlags.Name)

	_, err = client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}
