// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package available

import (
	"context"
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

type ListOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	AllNamespaces  bool

	Name string

	Summary bool
	Wide    bool

	positionalNameArg bool
}

func NewListOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger, positionalNameArg bool) *ListOptions {
	return &ListOptions{ui: ui, depsFactory: depsFactory, logger: logger, positionalNameArg: positionalNameArg}
}

func NewListCmd(o *ListOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	var examples cmdcore.Examples
	examples = append(examples,
		cmdcore.Example{"List packages available on the cluster",
			[]string{"package", "available", "list"},
		},
		cmdcore.Example{"List packages available on the cluster with their short descriptions",
			[]string{"package", "available", "list", "--wide"},
		},
		cmdcore.Example{"List all available package versions with release dates",
			[]string{"package", "available", "list", "--summary=false"},
		},
		cmdcore.Example{"List packages available in all namespaces",
			[]string{"package", "available", "list", "-A"},
		},
		cmdcore.Example{"List all available versions of a package",
			[]string{"package", "available", "list", "-p", "cert-manager.community.tanzu.vmware.com"},
		})

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List available packages in a namespace",
		RunE:    func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: examples.Description("kctrl", "-p", o.positionalNameArg),
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List available packages in all namespaces")

	cmd.Flags().BoolVar(&o.Summary, "summary", true, "Show summarized list of packages")
	if !o.positionalNameArg {
		cmd.Flags().StringVarP(&o.Name, "package", "p", "", "List all available versions of package")
	}

	cmd.Flags().BoolVar(&o.Wide, "wide", false, "Show additional info")

	return cmd
}

func (o *ListOptions) Run(args []string) error {
	if o.positionalNameArg && len(args) > 0 {
		o.Name = args[0]
	}

	if o.Summary && o.Name == "" {
		return o.listPackageMetadatas()
	}
	return o.listPackages()
}

func (o *ListOptions) listPackageMetadatas() error {
	tableTitle := fmt.Sprintf("Available summarized packages in namespace '%s'", o.NamespaceFlags.Name)
	nsHeader := uitable.NewHeader("Namespace")
	nsHeader.Hidden = true

	if o.AllNamespaces {
		o.NamespaceFlags.Name = ""
		tableTitle = "Available summarized packages in all namespaces"
		nsHeader.Hidden = false
	}

	client, err := o.depsFactory.PackageClient()
	if err != nil {
		return err
	}

	pkgmList, err := client.DataV1alpha1().PackageMetadatas(
		o.NamespaceFlags.Name).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	shortDesc := uitable.NewHeader("Short description")
	shortDesc.Hidden = !o.Wide

	table := uitable.Table{
		Title: tableTitle,

		Header: []uitable.Header{
			nsHeader,
			uitable.NewHeader("Name"),
			uitable.NewHeader("Display name"),
			shortDesc,
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	for _, pkgm := range pkgmList.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			cmdcore.NewValueNamespace(pkgm.Namespace),
			uitable.NewValueString(pkgm.Name),
			uitable.NewValueString(pkgm.Spec.DisplayName),
			uitable.NewValueString(pkgm.Spec.ShortDescription),
		})
	}

	o.ui.PrintTable(table)

	return err
}

func (o *ListOptions) listPackages() error {
	tableTitle := fmt.Sprintf("Available packages in namespace '%s'", o.NamespaceFlags.Name)
	nsHeader := uitable.NewHeader("Namespace")
	nsHeader.Hidden = true

	if o.AllNamespaces {
		o.NamespaceFlags.Name = ""
		tableTitle = "Available packages in all namespaces"
		nsHeader.Hidden = false
	}

	client, err := o.depsFactory.PackageClient()
	if err != nil {
		return err
	}

	listOpts := metav1.ListOptions{}
	if len(o.Name) > 0 {
		listOpts.FieldSelector = fields.Set{"spec.refName": o.Name}.String()
	}

	pkgList, err := client.DataV1alpha1().Packages(
		o.NamespaceFlags.Name).List(context.Background(), listOpts)
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title: tableTitle,

		Header: []uitable.Header{
			nsHeader,
			uitable.NewHeader("Name"),
			uitable.NewHeader("Version"),
			uitable.NewHeader("Released at"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
			{Column: 2, Asc: true},
		},
	}

	for _, pkg := range pkgList.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			cmdcore.NewValueNamespace(pkg.Namespace),
			uitable.NewValueString(pkg.Spec.RefName),
			cmdcore.NewValueSemver(pkg.Spec.Version),
			uitable.NewValueString(pkg.Spec.ReleasedAt.String()),
		})
	}

	o.ui.PrintTable(table)

	return err
}
