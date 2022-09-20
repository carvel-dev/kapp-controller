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

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewStatusOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *StatusOptions {
	return &StatusOptions{ui: ui, depsFactory: depsFactory, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewStatusCmd(o *StatusOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"s"},
		Short:   "View status of app created by package install",
		RunE:    func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Check status of package install",
				[]string{"package", "installed", "status", "-i", "cert-man"},
			},
		}.Description("-i", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations: map[string]string{cmdapp.TTYByDefaultKey: "",
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value},
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "package-install", "i", "", "Set installed package name (required)")
	} else {
		cmd.Use = "status INSTALLED_PACKAGE_NAME"
	}

	return cmd
}

func (o *StatusOptions) Run(args []string) error {
	if o.pkgCmdTreeOpts.PositionalArgs {
		o.Name = args[0]
	}

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
