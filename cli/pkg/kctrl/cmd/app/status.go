// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppStage string

const (
	fetchStage    AppStage = "fetch"
	templateStage AppStage = "template"
	deployStage   AppStage = "deploy"
	reconciled    AppStage = "reconciled"
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

	if o.IgnoreNotExists && !o.Follow {
		return fmt.Errorf("'--ignore-not-exists' can only be used while following")
	}

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	app, err := client.KappctrlV1alpha1().Apps(o.NamespaceFlags.Name).Get(context.Background(), o.Name, metav1.GetOptions{})
	if err != nil {
		if !(errors.IsNotFound(err) && o.IgnoreNotExists) {
			return err
		}
		o.ui.PrintLinef("AppCR '' in namespace '' does not exist...")
	}

	appWatcher := NewAppWatcher(o.NamespaceFlags.Name, o.Name, o.Follow, o.IgnoreNotExists, o.ui, client)

	if o.Follow {
		err = appWatcher.FollowApp(app)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = appWatcher.PrintTillCurrent(app.Status)
	if err != nil && !o.Follow {
		return err
	}

	return nil
}
