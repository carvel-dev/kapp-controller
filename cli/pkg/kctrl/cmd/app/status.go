// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
)

type StatusOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	IgnoreNotExists bool
}

func NewStatusOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *StatusOptions {
	return &StatusOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewStatusCmd(o *StatusOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"s"},
		Short:   "View status of app",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
		Annotations: map[string]string{TTYByDefaultKey: "",
			cmdcore.AppManagementCommandsHelpGroup.Key: cmdcore.AppManagementCommandsHelpGroup.Value},
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "app", "a", "", "Set app name (required)")
	cmd.Flags().BoolVar(&o.IgnoreNotExists, "ignore-not-exists", false, "Keep following app if it does not exist")

	return cmd
}

func (o *StatusOptions) Run() error {
	if len(o.Name) == 0 {
		return fmt.Errorf("Expected app name to be non empty")
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	appWatcher := NewAppTailer(o.NamespaceFlags.Name, o.Name, o.ui, client, AppTailerOpts{
		IgnoreNotExists:   o.IgnoreNotExists,
		PrintMetadata:     true,
		PrintCurrentState: true,
	})

	err = appWatcher.TailAppStatus()
	if err != nil {
		return err
	}

	return nil
}
