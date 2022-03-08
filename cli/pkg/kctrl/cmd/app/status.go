// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"

	"github.com/cppforlife/color"
	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatusOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	Follow          bool
	IgnoreNotExists bool
}

func NewStatusOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *StatusOptions {
	return &StatusOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewStatusCmd(o *StatusOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"s"},
		Short:   "View status of App CR",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "app", "a", "", "Set App CR name (required)")
	cmd.Flags().BoolVarP(&o.Follow, "follow", "f", false, "Follow changes in the App CRs status")
	cmd.Flags().BoolVar(&o.IgnoreNotExists, "ignore-not-exists", false, "Keep following AppCR if it does not exist")

	return cmd
}

func (o *StatusOptions) Run() error {
	if len(o.Name) == 0 {
		return fmt.Errorf("Expected App CR name to be non empty")
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	app, err := client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if o.Follow {
		o.ui.PrintLinef("Unimplemented")
		return nil
	}

	o.printStatus(app.Status)

	return nil
}

func (o *StatusOptions) printStatus(status kcv1alpha1.AppStatus) {
	if status.Fetch != nil {
		o.printHeader("Fetch")
		o.printStageMetadata(&status.Fetch.StartedAt, &status.Fetch.UpdatedAt)
		if status.Fetch.ExitCode != 0 && (status.Fetch.StartedAt.Before(&status.Fetch.UpdatedAt)) {
			o.ui.PrintBlock([]byte(color.RedString(status.Fetch.Error)))
			return
		}
		o.ui.PrintBlock([]byte(status.Fetch.Stdout))
		if status.Fetch.StartedAt.After(status.Fetch.UpdatedAt.Time) {
			o.printOngoing()
			return
		}
		o.printFinished()
	}

	if status.Template != nil {
		o.printHeader("Template")
		o.printStageMetadata(nil, &status.Template.UpdatedAt)
		if status.Template.ExitCode != 0 && (status.Fetch.StartedAt.Before(&status.Template.UpdatedAt)) {
			o.ui.PrintBlock([]byte(color.RedString(status.Template.Error)))
			return
		}
		if status.Fetch.StartedAt.After(status.Template.UpdatedAt.Time) {
			o.printOngoing()
			return
		}
		o.printFinished()
	}

	if status.Deploy != nil {
		o.printHeader("Deploy")
		o.printStageMetadata(&status.Deploy.StartedAt, &status.Deploy.UpdatedAt)
		if status.Deploy.ExitCode != 0 && (status.Deploy.StartedAt.Before(&status.Deploy.UpdatedAt)) {
			o.ui.PrintBlock([]byte(color.RedString(status.Fetch.Error)))
			return
		}
		o.ui.PrintBlock([]byte(status.Deploy.Stdout))
		if o.hasReconciled(status) {
			o.printFinished()
			return
		}
		o.printOngoing()
	}
}

func (o *StatusOptions) printHeader(header string) {
	o.ui.PrintLinef(color.New(color.Bold).Sprintf("-------------------%s-------------------", header))
}

func (o *StatusOptions) printStageMetadata(startedAt *metav1.Time, updatedAt *metav1.Time) {
	startedAtHeader := uitable.NewHeader("Started At")
	startedAtHeader.Hidden = (startedAt == nil)

	rows := []uitable.Value{
		nil,
		uitable.NewValueTime(updatedAt.Time),
	}

	if startedAt != nil {
		rows = []uitable.Value{
			uitable.NewValueTime(startedAt.Time),
			uitable.NewValueTime(updatedAt.Time),
		}
	}

	table := uitable.Table{
		Transpose: true,

		Header: []uitable.Header{
			startedAtHeader,
			uitable.NewHeader("Updated At"),
		},

		Rows: [][]uitable.Value{rows},
	}

	o.ui.PrintTable(table)
}

func (o *StatusOptions) hasReconciled(status kcv1alpha1.AppStatus) bool {
	for _, condition := range status.Conditions {
		if condition.Type == kcv1alpha1.ReconcileSucceeded && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func (o *StatusOptions) printFinished() {
	o.ui.PrintLinef(color.GreenString("Finished"))
}

func (o *StatusOptions) printOngoing() {
	o.ui.PrintLinef(color.YellowString("Ongoing"))
}
