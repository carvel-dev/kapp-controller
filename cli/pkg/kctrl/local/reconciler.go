// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package local

import (
	"context"
	"fmt"
	gourl "net/url"
	"os"
	"time"

	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	fakedpkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned/fake"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/app"
	fakekc "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	kcconfig "github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/metrics"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/packageinstall"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	depsFactory cmdcore.DepsFactory
	cmdRunner   exec.CmdRunner
	logger      logger.Logger
}

func NewReconciler(depsFactory cmdcore.DepsFactory,
	cmdRunner exec.CmdRunner, logger logger.Logger) *Reconciler {

	return &Reconciler{depsFactory, cmdRunner, logger}
}

type ReconcileOpts struct {
	Local     bool
	KbldBuild bool
	Delete    bool
	Debug     bool

	BeforeAppReconcile func(kcv1alpha1.App, *fakekc.Clientset) error
	AfterAppReconcile  func(kcv1alpha1.App, *fakekc.Clientset) error
}

func (o *Reconciler) Reconcile(configs Configs, opts ReconcileOpts) error {
	var objs []runtime.Object
	var appRes kcv1alpha1.App
	var primaryAnns map[string]string

	if len(configs.Apps) > 0 {
		appRes = configs.Apps[0]
		primaryAnns = appRes.Annotations
		if opts.Delete {
			appRes.DeletionTimestamp = &metav1.Time{time.Now()}
		}
		objs = append(objs, &appRes)
	}

	if len(configs.PkgInstalls) > 0 {
		pkgiRes := configs.PkgInstalls[0]
		primaryAnns = pkgiRes.Annotations
		// TODO delete does not delete because App CR does not exist in memory
		if opts.Delete {
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

	err = o.hackyConfigureKubernetesDst(coreClient)
	if err != nil {
		return err
	}

	minCoreClient := &MinCoreClient{
		client:          coreClient,
		localSecrets:    &localSecrets{configs.Secrets},
		localConfigMaps: configs.ConfigMaps,
	}
	kcClient := fakekc.NewSimpleClientset(objs...)
	dpkgClient := fakedpkg.NewSimpleClientset(configs.PkgsAsObjects()...)

	var vendirConfigHook func(conf vendirconf.Config) vendirconf.Config
	if opts.Local {
		vendirConf, err := newLocalVendirConf(primaryAnns)
		if err != nil {
			return fmt.Errorf("Calculating local vendir changes: %s", err)
		}
		vendirConfigHook = vendirConf.Adjust
	}

	appReconciler, pkgiReconciler := o.newReconcilers(
		minCoreClient, kcClient, dpkgClient, vendirConfigHook, opts)

	if opts.BeforeAppReconcile != nil {
		err := opts.BeforeAppReconcile(appRes, kcClient)
		if err != nil {
			return err
		}
	}

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

	if opts.AfterAppReconcile != nil {
		err := opts.AfterAppReconcile(appRes, kcClient)
		if err != nil {
			return err
		}
	}

	return reconcileErr
}

// hackyConfigureKubernetesDst configures environment variables for kapp.
// This would not be necessary if kapp was using default kubeconfig; however,
// right now kapp will use configuration based on configured serviceAccount within
// PackageInstall or App CR. However, we still need to configure it to know where to connect.
func (o *Reconciler) hackyConfigureKubernetesDst(coreClient kubernetes.Interface) error {
	host, err := o.depsFactory.RESTHost()
	if err != nil {
		return fmt.Errorf("Getting host: %s", err)
	}
	hostURL, err := gourl.Parse(host)
	if err != nil {
		return fmt.Errorf("Parsing host: %s", err)
	}
	os.Setenv("KUBERNETES_SERVICE_HOST", hostURL.Hostname())
	if hostURL.Port() == "" {
		os.Setenv("KUBERNETES_SERVICE_PORT", "443")
	} else {
		os.Setenv("KUBERNETES_SERVICE_PORT", hostURL.Port())
	}

	cm, err := coreClient.CoreV1().ConfigMaps("kube-public").Get(context.TODO(), "kube-root-ca.crt", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Fetching kube-root-ca.crt: %s", err)
	}
	// Used during fetching of service accounts in kapp-controller
	os.Setenv("KAPPCTRL_KUBERNETES_CA_DATA", cm.Data["ca.crt"])

	return nil
}

func (o *Reconciler) newReconcilers(
	coreClient kubernetes.Interface, kcClient *fakekc.Clientset, pkgClient *fakedpkg.Clientset,
	vendirConfigHook func(vendirconf.Config) vendirconf.Config, opts ReconcileOpts) (*app.Reconciler, *packageinstall.Reconciler) {

	runLog := logf.Log.WithName("deploy")
	if opts.Debug {
		// Only set logger in debug; logs go nowhere by default
		logf.SetLogger(zap.New(zap.UseDevMode(false)))
	}

	kcConfig := &kcconfig.Config{}

	appMetrics := metrics.NewAppMetrics()
	appMetrics.RegisterAllMetrics()

	refTracker := reftracker.NewAppRefTracker()
	updateStatusTracker := reftracker.NewAppUpdateStatus()

	appFactory := app.CRDAppFactory{
		CoreClient:       coreClient,
		AppClient:        kcClient,
		KcConfig:         kcConfig,
		AppMetrics:       appMetrics,
		VendirConfigHook: vendirConfigHook,
		KbldAllowBuild:   opts.KbldBuild, // only for CLI mode
		CmdRunner:        o.cmdRunner,
	}
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
