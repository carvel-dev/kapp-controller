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
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	AllNamespaces  bool
}

func NewListOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *ListOptions {
	return &ListOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewListCmd(o *ListOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	var examples cmdcore.Examples
	examples = append(examples,
		cmdcore.Example{"List package repositories",
			[]string{"package", "repository", "list"},
		},
		cmdcore.Example{"List package repositories in all namespaces",
			[]string{"package", "repository", "list", "A"}})

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List package repositories in a namespace",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
		Example: examples.Description("kctrl", "", false),
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List repositories in all namespaces")
	return cmd
}

func (o *ListOptions) Run() error {
	tableTitle := fmt.Sprintf("Repositories in namespace '%s'", o.NamespaceFlags.Name)
	nsHeader := uitable.NewHeader("Namespace")
	nsHeader.Hidden = true

	if o.AllNamespaces {
		o.NamespaceFlags.Name = ""
		tableTitle = "Repositories in all namespaces"
		nsHeader.Hidden = false
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	pkgrList, err := client.PackagingV1alpha1().PackageRepositories(
		o.NamespaceFlags.Name).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title: tableTitle,

		Header: []uitable.Header{
			nsHeader,
			uitable.NewHeader("Name"),
			uitable.NewHeader("Source"),
			uitable.NewHeader("Status"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	for _, pkgr := range pkgrList.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			cmdcore.NewValueNamespace(pkgr.Namespace),
			uitable.NewValueString(pkgr.Name),
			cmdcore.NewValueTruncated(NewSourceValue(pkgr), 60),
			cmdcore.NewValueTruncated(uitable.NewValueString(pkgr.Status.FriendlyDescription), 40),
		})
	}

	o.ui.PrintTable(table)

	return nil
}

// NewSourceValue returns a string summarizing spec.fetch for humans
// TODO should we place this into kapp-controller and expose as status field?
func NewSourceValue(pkgr v1alpha1.PackageRepository) uitable.Value {
	source := "(unknown)"

	if pkgr.Spec.Fetch != nil {
		switch {
		case pkgr.Spec.Fetch.ImgpkgBundle != nil:
			source = "(imgpkg) " + pkgr.Spec.Fetch.ImgpkgBundle.Image
			if pkgr.Spec.Fetch.ImgpkgBundle.TagSelection != nil && pkgr.Spec.Fetch.ImgpkgBundle.TagSelection.Semver != nil {
				source += fmt.Sprintf(" (%s)", pkgr.Spec.Fetch.ImgpkgBundle.TagSelection.Semver.Constraints)
			}
		default:
			// stays unknown
		}
	}

	return uitable.NewValueString(source)
}
