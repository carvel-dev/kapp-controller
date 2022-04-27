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

type DeleteOptions struct {
	ui          ui.UI
	statusUI    cmdcore.StatusLoggingUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	WaitFlags cmdcore.WaitFlags

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	NoOp bool
}

func NewDeleteOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *DeleteOptions {
	return &DeleteOptions{ui: ui, statusUI: cmdcore.NewStatusLoggingUI(ui), depsFactory: depsFactory, logger: logger}
}

func NewDeleteCmd(o *DeleteOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete app",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "app", "a", "", "Set app name (required)")
	cmd.Flags().BoolVar(&o.NoOp, "noop", false, "Ignore resources created by the app and delete the custom resource itself")
	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   5 * time.Minute,
	})

	return cmd
}

func (o *DeleteOptions) Run() error {
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
		o.ui.BeginLinef("App '%s' is owned by '%s'\n(The app will be created again when the package installation reconciles)\n", o.Name, fmt.Sprintf("%s/%s", app.OwnerReferences[0].Kind, app.OwnerReferences[0].Name))
	}

	o.ui.BeginLinef("Deleting app '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
	err = o.ui.AskForConfirmation()
	if err != nil {
		return err
	}

	if o.NoOp {
		err = o.patchNoopDelete(client)
		if err != nil {
			return err
		}
	}

	err = client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Delete(context.Background(), o.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	if o.WaitFlags.Enabled {
		err = o.waitForAppDeletion(client)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *DeleteOptions) patchNoopDelete(client kcclient.Interface) error {
	noopDeletePatch := []map[string]interface{}{
		{
			"op":    "add",
			"path":  "/spec/noopDelete",
			"value": true,
		},
	}

	patchJSON, err := json.Marshal(noopDeletePatch)
	if err != nil {
		return err
	}

	o.statusUI.PrintMessagef("Ignoring associated resources for app '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)

	_, err = client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (o *DeleteOptions) waitForAppDeletion(client kcclient.Interface) error {
	o.statusUI.PrintMessagef("Waiting for app deletion for '%s'", o.Name)
	appWatcher := NewAppTailer(o.NamespaceFlags.Name, o.Name, o.ui, client, AppTailerOpts{})
	err := appWatcher.TailAppStatus()
	if err != nil {
		return err
	}

	return nil
}
