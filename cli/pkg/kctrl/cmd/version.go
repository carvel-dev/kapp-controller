// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/version"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kappControllerNamespace  = "kapp-controller"
	kappControllerDeployment = "kapp-controller"
	kappControllerVersionAnn = "kapp-controller.carvel.dev/version"
)

type VersionOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory

	controllerVersion bool
}

func NewVersionOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *VersionOptions {
	return &VersionOptions{ui: ui, depsFactory: depsFactory}
}

func NewVersionCmd(o *VersionOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print client version",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	cmd.Flags().BoolVar(&o.controllerVersion, "controller", false, "Get the version of kapp-controller deployed on the cluster")

	return cmd
}

func (o *VersionOptions) Run() error {
	if o.controllerVersion {
		err := o.showControllerVersion()
		if err != nil {
			return err
		}
		return nil
	}

	o.ui.PrintBlock([]byte(fmt.Sprintf("kctrl version %s\n", version.Version)))

	return nil
}

func (o *VersionOptions) showControllerVersion() error {
	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	controllerDeployment, err := coreClient.AppsV1().Deployments(kappControllerNamespace).Get(context.Background(), kappControllerDeployment, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("kapp-controller not installed on cluster")
		}

		return err
	}

	o.ui.PrintBlock([]byte(fmt.Sprintf("kapp-controller version %s\n", controllerDeployment.GetAnnotations()[kappControllerVersionAnn])))

	return nil
}
