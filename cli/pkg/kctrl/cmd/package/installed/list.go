// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

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

type ListOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	AllNamespaces  bool

	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewListOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *ListOptions {
	return &ListOptions{ui: ui, depsFactory: depsFactory, logger: logger, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewListCmd(o *ListOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List installed packages in a namespace",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
		Example: cmdcore.Examples{
			cmdcore.Example{"List installed packages",
				[]string{"package", "installed", "list"},
			},
			cmdcore.Example{"List installed packages in all namespaces",
				[]string{"package", "installed", "list", "A"}},
		}.Description("", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations:  map[string]string{"table": ""},
	}
	o.NamespaceFlags.SetWithPackageCommandTreeOpts(cmd, flagsFactory, o.pkgCmdTreeOpts)
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List installed packages in all namespaces")
	return cmd
}

func (o *ListOptions) Run() error {
	tableTitle := fmt.Sprintf("Installed packages in namespace '%s'", o.NamespaceFlags.Name)
	nsHeader := uitable.NewHeader("Namespace")
	nsHeader.Hidden = true

	if o.AllNamespaces {
		o.NamespaceFlags.Name = ""
		tableTitle = "Installed packages in all namespaces"
		nsHeader.Hidden = false
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	pkgiList, err := client.PackagingV1alpha1().PackageInstalls(
		o.NamespaceFlags.Name).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title: tableTitle,

		Header: []uitable.Header{
			nsHeader,
			uitable.NewHeader("Name"),
			uitable.NewHeader("Package Name"),
			uitable.NewHeader("Package Version"),
			uitable.NewHeader("Status"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	for _, pkgi := range pkgiList.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			cmdcore.NewValueNamespace(pkgi.Namespace),
			uitable.NewValueString(pkgi.Name),
			uitable.NewValueString(pkgi.Spec.PackageRef.RefName),
			uitable.NewValueString(pkgi.Status.Version),
			uitable.NewValueString(packageInstallStatus(&pkgi)),
		})
	}

	o.ui.PrintTable(table)

	return nil
}
