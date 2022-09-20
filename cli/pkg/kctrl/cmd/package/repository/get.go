// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"context"
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewGetOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *GetOptions {
	return &GetOptions{ui: ui, depsFactory: depsFactory, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewGetCmd(o *GetOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get details for a package repository",
		Args:    cobra.ExactArgs(1),
		RunE:    func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Get details for a package repository",
				[]string{"package", "repository", "get", "-r", "tce"}},
		}.Description("-r", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations: map[string]string{"table": "",
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value},
	}
	o.NamespaceFlags.SetWithPackageCommandTreeOpts(cmd, flagsFactory, o.pkgCmdTreeOpts)

	if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "repository", "r", "", "Set package repository name (required)")
	} else {
		cmd.Use = "get REPOSITORY_NAME"
		cmd.Args = cobra.ExactArgs(1)
	}

	return cmd
}

func (o *GetOptions) Run(args []string) error {
	if o.pkgCmdTreeOpts.PositionalArgs {
		o.Name = args[0]
	}

	if len(o.Name) == 0 {
		return fmt.Errorf("Expected package repository name to be non-empty")
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	pkgr, err := client.PackagingV1alpha1().PackageRepositories(
		o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	table := uitable.Table{
		Transpose: true,

		Header: []uitable.Header{
			uitable.NewHeader("Namespace"),
			uitable.NewHeader("Name"),
			uitable.NewHeader("Source"),
			uitable.NewHeader("Status"),
			uitable.NewHeader("Conditions"),
			uitable.NewHeader("Useful error message"),
		},

		Rows: [][]uitable.Value{{
			uitable.NewValueString(pkgr.Namespace),
			uitable.NewValueString(pkgr.Name),
			NewSourceValue(*pkgr),
			uitable.NewValueString(pkgr.Status.FriendlyDescription),
			uitable.NewValueInterface(pkgr.Status.Conditions),
			uitable.NewValueString(pkgr.Status.UsefulErrorMessage),
		}},
	}

	o.ui.PrintTable(table)

	return nil
}
