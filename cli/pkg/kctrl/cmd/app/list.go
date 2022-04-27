// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"
	"time"

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
		Short:   "List Apps",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List apps in all namespaces")

	return cmd
}

func (o *ListOptions) Run() error {
	tableTitle := fmt.Sprintf("Available apps in namespace '%s'", o.NamespaceFlags.Name)
	nsHeader := uitable.NewHeader("Namespace")
	nsHeader.Hidden = true

	if o.AllNamespaces {
		o.NamespaceFlags.Name = ""
		tableTitle = "Available apps in all namespaces"
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
			uitable.NewHeader("Status"),
			uitable.NewHeader("Since Deploy"),
			uitable.NewHeader("Owner"),
			uitable.NewHeader("Age"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	for _, app := range appList.Items {
		sinceDeployAge := cmdcore.NewValueAge(time.Time{})
		if app.Status.Deploy != nil {
			sinceDeployAge = cmdcore.NewValueAge(app.Status.Deploy.UpdatedAt.Time)
		}

		table.Rows = append(table.Rows, []uitable.Value{
			cmdcore.NewValueNamespace(app.Namespace),
			uitable.NewValueString(app.Name),
			uitable.NewValueString(app.Status.FriendlyDescription),
			sinceDeployAge,
			uitable.NewValueString(o.owner(app.OwnerReferences)),
			cmdcore.NewValueAge(app.CreationTimestamp.Time),
		})
	}

	o.ui.PrintTable(table)

	return nil
}

func (o *ListOptions) owner(references []metav1.OwnerReference) string {
	if len(references) > 0 {
		return fmt.Sprintf("%s/%s", references[0].Kind, references[0].Name)
	}
	return ""
}
