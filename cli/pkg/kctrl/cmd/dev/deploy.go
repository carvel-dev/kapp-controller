// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package dev

import (
	"context"
	"fmt"
	gourl "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

	Files     []string
	Local     bool
	KbldBuild bool
	Delete    bool
	Debug     bool
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

	cmd.Flags().BoolVarP(&o.Local, "local", "l", false, "Use local fetch source")
	cmd.Flags().BoolVarP(&o.KbldBuild, "kbld-build", "b", false, "Allow kbld build")
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
	var primaryAnns map[string]string

	if len(configs.Apps) > 0 {
		appRes = configs.Apps[0]
		primaryAnns = appRes.Annotations
		if o.Delete {
			appRes.DeletionTimestamp = &metav1.Time{time.Now()}
		}
		objs = append(objs, &appRes)
	}

	if len(configs.PkgInstalls) > 0 {
		pkgiRes := configs.PkgInstalls[0]
		primaryAnns = pkgiRes.Annotations
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
	if o.Local {
		vendirConf, err := newLocalVendirConf(primaryAnns)
		if err != nil {
			return fmt.Errorf("Calculating local vendir changes: %s", err)
		}
		vendirConfigHook = vendirConf.Adjust
	}

	appReconciler, pkgiReconciler := o.newReconcilers(
		minCoreClient, kcClient, dpkgClient, vendirConfigHook)

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

// hackyConfigureKubernetesDst configures environment variables for kapp.
// This would not be necessary if kapp was using default kubeconfig; however,
// right now kapp will use configuration based on configured serviceAccount within
// PackageInstall or App CR. However, we still need to configure it to know where to connect.
func (o *DeployOptions) hackyConfigureKubernetesDst(coreClient kubernetes.Interface) error {
	host, err := o.depsFactory.RESTHost()
	if err != nil {
		return fmt.Errorf("Getting host: %s", err)
	}
	hostURL, err := gourl.Parse(host)
	if err != nil {
		return fmt.Errorf("Parsing host: %s", err)
	}
	os.Setenv("KUBERNETES_SERVICE_HOST", hostURL.Hostname())
	os.Setenv("KUBERNETES_SERVICE_PORT", hostURL.Port())

	cm, err := coreClient.CoreV1().ConfigMaps("kube-public").Get(context.TODO(), "kube-root-ca.crt", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Fetching kube-root-ca.crt: %s", err)
	}
	// Used during fetching of service accounts in kapp-controller
	os.Setenv("KAPPCTRL_KUBERNETES_CA_DATA", cm.Data["ca.crt"])

	return nil
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
	coreClient kubernetes.Interface, kcClient *fakekc.Clientset, pkgClient *fakedpkg.Clientset,
	vendirConfigHook func(vendirconf.Config) vendirconf.Config) (*app.Reconciler, *packageinstall.Reconciler) {

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

	appFactory := app.CRDAppFactory{
		CoreClient:       coreClient,
		AppClient:        kcClient,
		KcConfig:         kcConfig,
		AppMetrics:       appMetrics,
		VendirConfigHook: vendirConfigHook,
		KbldAllowBuild:   o.KbldBuild, // only for CLI mode
		CmdRunner:        NewDetailedCmdRunner(os.Stdout, o.Debug),
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

type localVendirConf struct {
	// Indexed by numeric fetch index
	localPaths map[int]string
}

func newLocalVendirConf(resAnnotations map[string]string) (localVendirConf, error) {
	cwdPath, err := os.Getwd()
	if err != nil {
		return localVendirConf{}, err
	}

	const (
		prefix = "kctrl.carvel.dev/local-fetch-"
	)

	localPaths := map[int]string{}

	for key, val := range resAnnotations {
		if strings.HasPrefix(key, prefix) {
			fetchIdx, err := strconv.Atoi(strings.TrimPrefix(key, prefix))
			if err != nil {
				return localVendirConf{}, err
			}
			localPaths[fetchIdx] = filepath.Join(cwdPath, val)
		}
	}

	return localVendirConf{localPaths}, nil
}

func (c localVendirConf) Adjust(conf vendirconf.Config) vendirconf.Config {
	for fetchIdx, localPath := range c.localPaths {
		if fetchIdx >= len(conf.Directories) {
			// Ignore invalid indexes
			continue
		}
		conf.Directories[fetchIdx].Contents[0] = vendirconf.DirectoryContents{
			Path: conf.Directories[fetchIdx].Contents[0].Path,
			Directory: &vendirconf.DirectoryContentsDirectory{
				Path: localPath,
			},
		}
	}
	return conf
}
