// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageavailable

import (
	"context"
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	cmdcore "github.com/k14s/kapp/pkg/kapp/cmd/core"
	"github.com/k14s/kapp/pkg/kapp/logger"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	AllNamespaces  bool
	PackageName    string
}

func NewListOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *ListOptions {
	return &ListOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewListCmd(o *ListOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List available packages in a namespace",
		Args:    cobra.MaximumNArgs(1),
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List available packages")
	cmd.Flags().StringVarP(&o.PackageName, "package", "P", "", "List all available versions of package")
	return cmd
}

func (o *ListOptions) Run() error {
	var table uitable.Table
	var err error
	if o.PackageName != "" {
		table, err = listAvailablePackageVersions(o)
	} else {
		table, err = listAvailablePackages(o)
	}
	if err != nil {
		return err
	}

	o.ui.PrintTable(table)

	return nil
}

func listAvailablePackages(o *ListOptions) (uitable.Table, error) {
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
		return uitable.Table{}, err
	}

	pkgaList, err := client.DataV1alpha1().PackageMetadatas(
		o.NamespaceFlags.Name).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return uitable.Table{}, err
	}

	table := uitable.Table{
		Title:   tableTitle,
		Content: "Packages Available",

		Header: []uitable.Header{
			nsHeader,
			uitable.NewHeader("Name"),
			uitable.NewHeader("Display-Name"),
			uitable.NewHeader("Short-Description"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	for _, pkga := range pkgaList.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			cmdcore.NewValueNamespace(pkga.Namespace),
			uitable.NewValueString(pkga.Name),
			uitable.NewValueString(pkga.Spec.DisplayName),
			uitable.NewValueString(pkga.Spec.ShortDescription),
		})
	}
	return table, err
}

func listAvailablePackageVersions(o *ListOptions) (uitable.Table, error) {
	tableTitle := fmt.Sprintf("Available package versions for '%s' in namespace '%s'", o.PackageName, o.NamespaceFlags.Name)
	nsHeader := uitable.NewHeader("Namespace")
	nsHeader.Hidden = true

	if o.AllNamespaces {
		o.NamespaceFlags.Name = ""
		tableTitle = fmt.Sprintf("Available package versions for '%s' in all namespaces", o.PackageName)
		nsHeader.Hidden = false
	}

	client, err := o.depsFactory.PackageClient()
	if err != nil {
		return uitable.Table{}, err
	}

	fieldSelector := fmt.Sprintf("spec.refName=%s", o.PackageName)
	pkgaList, err := client.DataV1alpha1().Packages(
		o.NamespaceFlags.Name).List(context.Background(), metav1.ListOptions{FieldSelector: fieldSelector})
	if err != nil {
		return uitable.Table{}, err
	}

	table := uitable.Table{
		Title:   tableTitle,
		Content: "Package Versions Available",

		Header: []uitable.Header{
			nsHeader,
			uitable.NewHeader("Name"),
			uitable.NewHeader("Version"),
			uitable.NewHeader("Released-At"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	for _, pkga := range pkgaList.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			cmdcore.NewValueNamespace(pkga.Namespace),
			uitable.NewValueString(pkga.Spec.RefName),
			uitable.NewValueString(pkga.Spec.Version),
			uitable.NewValueString(pkga.Spec.ReleasedAt.String()),
		})
	}
	return table, err
}
