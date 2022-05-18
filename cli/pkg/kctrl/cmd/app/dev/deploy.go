// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package dev

import (
	"context"
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdapp "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/app"
	fakekc "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	kcconfig "github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/metrics"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	"k8s.io/apimachinery/pkg/types"
	// fakecore "k8s.io/client-go/kubernetes/fake"
	// corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/yaml"
)

type DeployOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	Name           string

	Files  []string
	Delete bool
	Debug  bool
}

func NewDeployOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *DeployOptions {
	return &DeployOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewDeployCmd(o *DeployOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy App CR",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.Name, "app", "a", "", "Set App CR name (required)")
	cmd.Flags().StringSliceVarP(&o.Files, "file", "f", nil, "Set App CR file (required)")

	cmd.Flags().BoolVar(&o.Delete, "delete", false, "Delete deployed app")
	cmd.Flags().BoolVar(&o.Debug, "debug", false, "Show kapp-controller logs")

	return cmd
}

func (o *DeployOptions) Run() error {
	appRes, secrets, configMaps, err := NewConfigFromFiles(o.Files)
	if err != nil {
		return fmt.Errorf("Reading App CR configuration files: %s", err)
	}

	if len(o.Name) > 0 {
		appRes.Name = o.Name
	}
	// Prefer namespace specified in the configuration
	if len(appRes.Namespace) == 0 {
		appRes.Namespace = o.NamespaceFlags.Name
	}
	if o.Delete {
		appRes.DeletionTimestamp = &metav1.Time{time.Now()}
	}

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return fmt.Errorf("Getting core client: %s", err)
	}
	minCoreClient := &MinCoreClient{
		client:          coreClient,
		localSecrets:    secrets,
		localConfigMaps: configMaps,
	}
	kcClient := fakekc.NewSimpleClientset(&appRes)
	reconciler := o.newReconciler(minCoreClient, kcClient)

	err = o.printAppRes(appRes, kcClient)
	if err != nil {
		return err
	}

	o.ui.PrintLinef("Reconciling in-memory app/%s (namespace: %s) ...", appRes.Name, appRes.Namespace)

	go func() {
		appWatcher := cmdapp.NewAppTailer(appRes.Namespace, appRes.Name, o.ui, kcClient, cmdapp.AppTailerOpts{})

		err := appWatcher.TailAppStatus()
		if err != nil {
			o.ui.PrintLinef("tailing error: %s", err)
		}
	}()

	// TODO is there a better way to deal with service accounts?
	// TODO do anything with reconcile result?
	_, reconcileErr := reconciler.Reconcile(context.TODO(), reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      appRes.Name,
			Namespace: appRes.Namespace,
		},
	})

	if o.Debug {
		err := o.printAppRes(appRes, kcClient)
		if err != nil {
			return err
		}
	}

	return reconcileErr
}

func (o *DeployOptions) printAppRes(app kcv1alpha1.App, kcClient *fakekc.Clientset) error {
	updatedApp, err := kcClient.KappctrlV1alpha1().Apps(app.Namespace).Get(context.Background(), app.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Fetching App CR: %s", err)
	}

	bs, err := yaml.Marshal(updatedApp)
	if err != nil {
		return fmt.Errorf("Marshaling App CR: %s", err)
	}

	o.ui.PrintBlock(bs)

	return nil
}

func (o *DeployOptions) newReconciler(
	coreClient kubernetes.Interface, kcClient *fakekc.Clientset) *app.Reconciler {

	runLog := logf.Log.WithName("deploy")
	if o.Debug {
		// Only set logger in debug; logs go nowhere by default
		logf.SetLogger(zap.New(zap.UseDevMode(false)))
	}

	kcConfig := &kcconfig.Config{}

	appMetrics := metrics.NewAppMetrics()
	appMetrics.RegisterAllMetrics()

	refTracker := reftracker.NewAppRefTracker()
	updateStatusTracker := reftracker.NewAppUpdateStatus()

	appFactory := app.CRDAppFactory{coreClient, kcClient, kcConfig, appMetrics}
	return app.NewReconciler(kcClient, runLog, appFactory, refTracker, updateStatusTracker)
}
