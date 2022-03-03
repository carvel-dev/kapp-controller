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
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

type DeleteOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	WaitFlags cmdcore.WaitFlags

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	IgnoreAssociatedResources bool
}

func NewDeleteOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *DeleteOptions {
	return &DeleteOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewDeleteCmd(o *DeleteOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete App CR",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "app", "a", "", "Set App CR name (required)")
	cmd.Flags().BoolVar(&o.IgnoreAssociatedResources, "ignore-associated-resources", false, "Ignore resources created by the AppCR and delete the custom resource itself")
	o.WaitFlags.Set(cmd, flagsFactory, &cmdcore.WaitFlagsOpts{
		AllowDisableWait: true,
		DefaultInterval:  1 * time.Second,
		DefaultTimeout:   5 * time.Minute,
	})

	return cmd
}

func (o *DeleteOptions) Run() error {
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
		return fmt.Errorf("App CR '%s' in namespace '%s' is owned by a PackageInstall.\n(Hint: Try using `kctrl package installed delete` to delete PackageInstall)",
			o.Name, o.NamespaceFlags.Name)
	}

	o.ui.PrintLinef("Deleting App CR '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
	err = o.ui.AskForConfirmation()
	if err != nil {
		return err
	}

	if o.IgnoreAssociatedResources {
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

	o.ui.PrintLinef("Ignoring associated resources for App CR '%s' in namespace '%s'...", o.Name, o.NamespaceFlags.Name)

	_, err = client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Patch(context.Background(), o.Name, types.JSONPatchType, patchJSON, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (o *DeleteOptions) waitForAppDeletion(client kcclient.Interface) error {
	o.ui.PrintLinef("Waiting for App CR deletion for '%s'", o.Name)
	msgsUI := cmdcore.NewDedupingMessagesUI(cmdcore.NewPlainMessagesUI(o.ui))
	description := getAppDescription(o.Name, o.NamespaceFlags.Name)

	err := wait.Poll(o.WaitFlags.CheckInterval, o.WaitFlags.Timeout, func() (bool, error) {
		app, err := client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				msgsUI.NotifySection("%s: DeletionSucceeded", description)
				return true, nil
			}
			return false, err
		}
		if app.Generation != app.Status.ObservedGeneration {
			// Should wait for generation to be observed before checking the reconciliation status so that we know we are checking the new spec
			return false, nil
		}
		status := app.Status.GenericStatus

		for _, cond := range status.Conditions {
			msgsUI.NotifySection("%s: %s", description, cond.Type)

			if cond.Type == kcv1alpha1.DeleteFailed && cond.Status == corev1.ConditionTrue {
				return false, fmt.Errorf("%s: Deleting: %s. %s", description, status.UsefulErrorMessage, status.FriendlyDescription)
			}
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("%s: Deleting: %s", description, err)
	}

	return nil
}
