// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"context"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kapp/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string
}

func NewGetOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *GetOptions {
	return &GetOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewGetCmd(o *GetOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get details for a package repository",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "repository", "r", "", "Set package repository name")
	return cmd
}

func (o *GetOptions) Run() error {
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
			uitable.NewHeader("Description"),
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
