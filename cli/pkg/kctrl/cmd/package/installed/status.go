// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	"context"
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdapp "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatusOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	IgnoreNotExists bool
}

func NewStatusOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *StatusOptions {
	return &StatusOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewStatusCmd(o *StatusOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"s"},
		Short:   "View status of app created by package install",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "package-install", "i", "", "Set package installname (required)")

	return cmd
}

func (o *StatusOptions) Run() error {
	if len(o.Name) == 0 {
		return fmt.Errorf("Expected package install name to be non empty")
	}

	// TODO: Should we assert that pakage install exists?

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	app, err := client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if app.OwnerReferences[0].Kind != "PackageInstall" || app.OwnerReferences[0].Name != o.Name {
		return fmt.Errorf("Could not find app associated with package install '%s' in namespace '%s'", o.Name, o.NamespaceFlags.Name)
	}

	appWatcher := cmdapp.NewAppTailer(o.NamespaceFlags.Name, o.Name, o.ui, client, cmdapp.AppTailerOpts{
		PrintMetadata:     true,
		PrintCurrentState: true,
	})

	err = appWatcher.TailAppStatus()
	if err != nil {
		return err
	}

	return nil
}
