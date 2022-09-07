// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package dev

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdapp "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	cmdlocal "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	fakekc "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type DevOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags

	Files     []string
	Local     bool
	KbldBuild bool
	Delete    bool
	Debug     bool
}

func NewDevOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *DevOptions {
	return &DevOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewCmd(o *DevOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Deploy App CRs and packaging CRs",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
		Annotations: map[string]string{
			cmdcore.DevHelpGroup.Key: cmdcore.DevHelpGroup.Value,
		},
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringSliceVarP(&o.Files, "file", "f", nil, "Set App CR file (required)")

	cmd.Flags().BoolVarP(&o.Local, "local", "l", false, "Use local fetch source")
	cmd.Flags().BoolVarP(&o.KbldBuild, "kbld-build", "b", false, "Allow kbld build")
	cmd.Flags().BoolVar(&o.Delete, "delete", false, "Delete deployed app")
	cmd.Flags().BoolVar(&o.Debug, "debug", false, "Show kapp-controller logs")

	return cmd
}

func (o *DevOptions) Run() error {
	configs, err := cmdlocal.NewConfigFromFiles(o.Files)
	if err != nil {
		return fmt.Errorf("Reading App CR configuration files: %s", err)
	}

	configs.ApplyNamespace(o.NamespaceFlags.Name)

	cmdRunner := cmdlocal.NewDetailedCmdRunner(os.Stdout, o.Debug)
	reconciler := cmdlocal.NewReconciler(o.depsFactory, cmdRunner, o.logger)

	reconcileErr := reconciler.Reconcile(configs, cmdlocal.ReconcileOpts{
		Local:     o.Local,
		KbldBuild: o.KbldBuild,
		Delete:    o.Delete,
		Debug:     o.Debug,

		BeforeAppReconcile: o.beforeAppReconcile,
		AfterAppReconcile:  o.afterAppReconcile,
	})

	// TODO app watcher needs a little time to run; should block ideally
	time.Sleep(100 * time.Millisecond)

	return reconcileErr
}

func (o *DevOptions) beforeAppReconcile(app kcv1alpha1.App, kcClient *fakekc.Clientset) error {
	err := o.printRs(app.ObjectMeta, kcClient)
	if err != nil {
		return err
	}

	o.ui.PrintLinef("Reconciling in-memory app/%s (namespace: %s) ...", app.Name, app.Namespace)

	go func() {
		appWatcher := cmdapp.NewAppTailer(app.Namespace, app.Name,
			o.ui, kcClient, cmdapp.AppTailerOpts{IgnoreNotExists: true})

		err := appWatcher.TailAppStatus()
		if err != nil {
			o.ui.PrintLinef("App tailing error: %s", err)
		}
	}()

	return nil
}

func (o *DevOptions) afterAppReconcile(app kcv1alpha1.App, kcClient *fakekc.Clientset) error {
	if o.Debug {
		return o.printRs(app.ObjectMeta, kcClient)
	}
	return nil
}

func (o *DevOptions) printRs(nsName metav1.ObjectMeta, kcClient *fakekc.Clientset) error {
	app, err := kcClient.KappctrlV1alpha1().Apps(nsName.Namespace).Get(context.Background(), nsName.Name, metav1.GetOptions{})
	if err == nil {
		bs, err := yaml.Marshal(app)
		if err != nil {
			return fmt.Errorf("Marshaling App CR: %s", err)
		}

		o.ui.PrintBlock(bs)
	}

	pkgi, err := kcClient.PackagingV1alpha1().PackageInstalls(nsName.Namespace).Get(context.Background(), nsName.Name, metav1.GetOptions{})
	if err == nil {
		bs, err := yaml.Marshal(pkgi)
		if err != nil {
			return fmt.Errorf("Marshaling PackageInstall CR: %s", err)
		}

		o.ui.PrintBlock(bs)
	}

	return nil
}
