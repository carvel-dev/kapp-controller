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
	fakedpkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned/fake"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/app"
	fakekc "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	kcconfig "github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/metrics"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/packageinstall"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	cmd.Flags().StringSliceVarP(&o.Files, "file", "f", nil, "Set App CR file (required)")

	cmd.Flags().BoolVar(&o.Delete, "delete", false, "Delete deployed app")
	cmd.Flags().BoolVar(&o.Debug, "debug", false, "Show kapp-controller logs")

	return cmd
}

func (o *DeployOptions) Run() error {
	configs, err := NewConfigFromFiles(o.Files)
	if err != nil {
		return fmt.Errorf("Reading App CR configuration files: %s", err)
	}

	configs.ApplyNamespace(o.NamespaceFlags.Name)

	var objs []runtime.Object
	var appRes kcv1alpha1.App

	if len(configs.Apps) > 0 {
		appRes = configs.Apps[0]
		if o.Delete {
			appRes.DeletionTimestamp = &metav1.Time{time.Now()}
		}
		objs = append(objs, &appRes)
	}

	if len(configs.PkgInstalls) > 0 {
		pkgiRes := configs.PkgInstalls[0]
		// TODO delete does not delete because App CR does not exist in memory
		if o.Delete {
			pkgiRes.DeletionTimestamp = &metav1.Time{time.Now()}
		}
		objs = append(objs, &pkgiRes)

		// Specifies underlying app resource
		appRes.Name = pkgiRes.Name
		appRes.Namespace = pkgiRes.Namespace
	}

	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return fmt.Errorf("Getting core client: %s", err)
	}
	minCoreClient := &MinCoreClient{
		client:          coreClient,
		localSecrets:    &localSecrets{configs.Secrets},
		localConfigMaps: configs.ConfigMaps,
	}
	kcClient := fakekc.NewSimpleClientset(objs...)
	dpkgClient := fakedpkg.NewSimpleClientset(configs.PkgsAsObjects()...)
	appReconciler, pkgiReconciler := o.newReconcilers(minCoreClient, kcClient, dpkgClient)

	err = o.printRs(appRes.ObjectMeta, kcClient)
	if err != nil {
		return err
	}

	o.ui.PrintLinef("Reconciling in-memory app/%s (namespace: %s) ...", appRes.Name, appRes.Namespace)

	go func() {
		appWatcher := cmdapp.NewAppTailer(appRes.Namespace, appRes.Name,
			o.ui, kcClient, cmdapp.AppTailerOpts{IgnoreNotExists: true})

		err := appWatcher.TailAppStatus()
		if err != nil {
			o.ui.PrintLinef("App tailing error: %s", err)
		}
	}()

	var reconcileErr error

	if len(configs.PkgInstalls) > 0 {
		_, reconcileErr = pkgiReconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      configs.PkgInstalls[0].Name,
				Namespace: configs.PkgInstalls[0].Namespace,
			},
		})
	}

	// TODO is there a better way to deal with service accounts?
	// TODO do anything with reconcile result?
	_, reconcileErr = appReconciler.Reconcile(context.TODO(), reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      appRes.Name,
			Namespace: appRes.Namespace,
		},
	})

	// One more time to get successful or failed status
	if len(configs.PkgInstalls) > 0 {
		_, reconcileErr = pkgiReconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      configs.PkgInstalls[0].Name,
				Namespace: configs.PkgInstalls[0].Namespace,
			},
		})
	}

	if o.Debug {
		err := o.printRs(appRes.ObjectMeta, kcClient)
		if err != nil {
			return err
		}
	}

	// TODO app watcher needs a little time to run; should block ideally
	time.Sleep(100 * time.Millisecond)

	return reconcileErr
}

func (o *DeployOptions) printRs(nsName metav1.ObjectMeta, kcClient *fakekc.Clientset) error {
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

func (o *DeployOptions) newReconcilers(
	coreClient kubernetes.Interface, kcClient *fakekc.Clientset,
	pkgClient *fakedpkg.Clientset) (*app.Reconciler, *packageinstall.Reconciler) {

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
	appReconciler := app.NewReconciler(kcClient, runLog.WithName("app"),
		appFactory, refTracker, updateStatusTracker)

	pkgiReconciler := packageinstall.NewReconciler(
		kcClient, pkgClient, coreClient,
		// TODO do not need this in the constructor of Reconciler
		(*packageinstall.PackageInstallVersionHandler)(nil),
		runLog.WithName("pkgi"),
	)

	return appReconciler, pkgiReconciler
}
