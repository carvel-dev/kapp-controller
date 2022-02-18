// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

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
}

func NewListOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *ListOptions {
	return &ListOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewListCmd(o *ListOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"g"},
		Short:   "List App CRs",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List App CRs in all namespaces")

	return cmd
}

func (o *ListOptions) Run() error {
	tableTitle := fmt.Sprintf("Available App CRs in namespace '%s'", o.NamespaceFlags.Name)
	nsHeader := uitable.NewHeader("Namespace")
	nsHeader.Hidden = true

	if o.AllNamespaces {
		o.NamespaceFlags.Name = ""
		tableTitle = "Available App CRs in all namespaces"
		nsHeader.Hidden = false
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	appList, err := client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title: tableTitle,

		Header: []uitable.Header{
			nsHeader,
			uitable.NewHeader("Name"),
			uitable.NewHeader("Description"),
			uitable.NewHeader("Since Deploy"),
			uitable.NewHeader("Age"),
			uitable.NewHeader("Owned"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	for _, app := range appList.Items {
		table.Rows = append(table.Rows, []uitable.Value{
			cmdcore.NewValueNamespace(app.Namespace),
			uitable.NewValueString(app.Name),
			uitable.NewValueString(app.Status.FriendlyDescription),
			cmdcore.NewValueAge(app.Status.Deploy.UpdatedAt.Time),
			cmdcore.NewValueAge(app.CreationTimestamp.Time),
			uitable.NewValueString(o.isOwned(app.OwnerReferences)),
		})
	}

	o.ui.PrintTable(table)

	return nil
}

func (o *ListOptions) isOwned(references []metav1.OwnerReference) string {
	if len(references) > 0 {
		return "true"
	}
	return "false"
}
