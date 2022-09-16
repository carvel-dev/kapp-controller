// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"net/http"         // Pprof related
	_ "net/http/pprof" // Pprof related
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver"
	pkgclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/app"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/componentinfo"
	kcconfig "github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/kubeconfig"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/metrics"
	pkginstall "github.com/vmware-tanzu/carvel-kapp-controller/pkg/packageinstall"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/pkgrepository"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/sidecarexec"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // Initialize gcp client auth plugin
	"k8s.io/component-base/cli/flag"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	PprofListenAddr       = "0.0.0.0:6060"
	kappctrlAPIPORTEnvKey = "KAPPCTRL_API_PORT"
)

type Options struct {
	Concurrency            int
	Namespace              string
	EnablePprof            bool
	APIRequestTimeout      time.Duration
	PackagingGloablNS      string
	MetricsBindAddress     string
	APIPriorityAndFairness bool
	StartAPIServer         bool
	TLSCipherSuites        string
}

// Based on https://github.com/kubernetes-sigs/controller-runtime/blob/8f633b179e1c704a6e40440b528252f147a3362a/examples/builtins/main.go
func Run(opts Options, runLog logr.Logger) error {
	runLog.Info("start controller")
	runLog.Info("setting up manager")

	restConfig := config.GetConfigOrDie()

	if opts.APIRequestTimeout != 0 {
		restConfig.Timeout = opts.APIRequestTimeout
	}

	mgr, err := manager.New(restConfig, manager.Options{Namespace: opts.Namespace,
		Scheme: kcconfig.Scheme, MetricsBindAddress: opts.MetricsBindAddress})
	if err != nil {
		return fmt.Errorf("Setting up overall controller manager: %s", err)
	}

	runLog.Info("setting up controller")

	coreClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("Building core client: %s", err)
	}

	kcClient, err := kcclient.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("Building kappctrl client: %s", err)
	}

	kcConfig, err := kcconfig.NewConfig(coreClient)
	if err != nil {
		return fmt.Errorf("getting kapp-controller config: %s", err)
	}

	pkgClient, err := pkgclient.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("Building packaging client: %s", err)
	}

	runLog.Info("setting up metrics")
	appMetrics := metrics.NewAppMetrics()
	appMetrics.RegisterAllMetrics()

	var server *apiserver.APIServer
	if opts.StartAPIServer {
		// assign bindPort to env var KAPPCTRL_API_PORT if available
		var bindPort int
		if apiPort, ok := os.LookupEnv(kappctrlAPIPORTEnvKey); ok {
			var err error
			if bindPort, err = strconv.Atoi(apiPort); err != nil {
				return fmt.Errorf("Reading %s env var (must be int): %s", kappctrlAPIPORTEnvKey, err)
			}
		} else {
			return fmt.Errorf("Expected to find %s env var", kappctrlAPIPORTEnvKey)
		}

		// to facilitate creation of many packages at once from a larger PKGR
		pkgRestConfig := config.GetConfigOrDie()
		pkgRestConfig.QPS = 60
		pkgRestConfig.Burst = 90
		pkgKcClient, err := kcclient.NewForConfig(pkgRestConfig)
		if err != nil {
			return fmt.Errorf("Building pkg kappctrl client: %s", err)
		}

		cSuites, err := parseTLSCipherSuites(opts.TLSCipherSuites)
		if err != nil {
			return err
		}

		server, err := apiserver.NewAPIServer(pkgRestConfig, coreClient, pkgKcClient, apiserver.NewAPIServerOpts{
			GlobalNamespace:              opts.PackagingGloablNS,
			BindPort:                     bindPort,
			EnableAPIPriorityAndFairness: opts.APIPriorityAndFairness,
			Logger:                       runLog.WithName("apiserver"),
			TLSCipherSuites:              cSuites,
		})
		if err != nil {
			return fmt.Errorf("Building API server: %s", err)
		}

		err = server.Run()
		if err != nil {
			return fmt.Errorf("Starting API server: %s", err)
		}
	}

	sidecarClient, err := sidecarexec.NewClient(exec.NewPlainCmdRunner())
	if err != nil {
		return fmt.Errorf("Starting RPC client: %s", err)
	}

	sidecarCmdExec := sidecarClient.CmdExec()

	{ // add controller for config
		reconciler := kcconfig.NewReconciler(
			coreClient, kcConfig, sidecarClient.OSConfig(), runLog.WithName("config"))

		ctrl, err := controller.New("config", mgr, controller.Options{
			Reconciler:              reconciler,
			MaxConcurrentReconciles: 1,
		})
		if err != nil {
			return fmt.Errorf("Setting up Config reconciler: %s", err)
		}

		ns := os.Getenv("KAPPCTRL_SYSTEM_NAMESPACE")
		if ns == "" {
			return fmt.Errorf("Cannot get kapp-controller namespace")
		}

		err = reconciler.AttachWatches(ctrl, ns)
		if err != nil {
			return fmt.Errorf("Setting up Config reconciler watches: %s", err)
		}

		// Reconcile once synchronously to ensure controller configuration
		// (e.g. proxy, CA certs) is applied to sidecar before any tool execution happens.
		_, err = reconciler.Reconcile(context.TODO(), reconcile.Request{})
		if err != nil {
			return fmt.Errorf("Reconcile config reconciler once: %s", err)
		}
	}

	refTracker := reftracker.NewAppRefTracker()
	updateStatusTracker := reftracker.NewAppUpdateStatus()

	// initialize cluster access once - it contains a service account token cache which should be only setup once.
	kubeconf := kubeconfig.NewKubeconfig(coreClient, runLog)
	compInfo := componentinfo.NewComponentInfo(coreClient, kubeconf, Version)

	{ // add controller for apps
		appFactory := app.CRDAppFactory{
			CoreClient: coreClient,
			AppClient:  kcClient,
			KcConfig:   kcConfig,
			AppMetrics: appMetrics,
			CmdRunner:  sidecarCmdExec,
			Kubeconf:   kubeconf,
			CompInfo:   compInfo,
		}
		reconciler := app.NewReconciler(kcClient, runLog.WithName("app"),
			appFactory, refTracker, updateStatusTracker, compInfo)

		ctrl, err := controller.New("app", mgr, controller.Options{
			Reconciler: NewUniqueReconciler(&ErrReconciler{
				delegate: reconciler,
				log:      runLog.WithName("er"),
			}),
			MaxConcurrentReconciles: opts.Concurrency,
		})
		if err != nil {
			return fmt.Errorf("Setting up Apps reconciler: %s", err)
		}

		err = reconciler.AttachWatches(ctrl)
		if err != nil {
			return fmt.Errorf("Setting up Apps reconciler watches: %s", err)
		}
	}

	{ // add controller for PackageInstall
		pkgToPkgInstallHandler := pkginstall.NewPackageInstallVersionHandler(
			kcClient, opts.PackagingGloablNS, runLog.WithName("handler"))

		reconciler := pkginstall.NewReconciler(kcClient, pkgClient, coreClient, pkgToPkgInstallHandler, runLog.WithName("pkgi"), compInfo)

		ctrl, err := controller.New("pkgi", mgr, controller.Options{
			Reconciler:              reconciler,
			MaxConcurrentReconciles: 1,
		})
		if err != nil {
			return fmt.Errorf("Setting up PackageInstalls reconciler: %s", err)
		}

		err = reconciler.AttachWatches(ctrl)
		if err != nil {
			return fmt.Errorf("Setting up PackageInstalls reconciler watches: %s", err)
		}
	}

	{ // add controller for pkgrepositories
		appFactory := pkgrepository.AppFactory{coreClient, kcClient, kcConfig, sidecarCmdExec, kubeconf}

		reconciler := pkgrepository.NewReconciler(kcClient, coreClient,
			runLog.WithName("pkgr"), appFactory, refTracker, updateStatusTracker)

		ctrl, err := controller.New("pkgr", mgr, controller.Options{
			Reconciler: reconciler,
			// TODO: Consider making this configurable for multiple PackageRepo reconciles
			MaxConcurrentReconciles: 1,
		})
		if err != nil {
			return fmt.Errorf("Setting up PackageRepositories reconciler: %s", err)
		}

		err = reconciler.AttachWatches(ctrl)
		if err != nil {
			return fmt.Errorf("Setting up PackageRepositories reconciler watches: %s", err)
		}
	}

	runLog.Info("starting manager")

	if opts.EnablePprof {
		runLog.Info("DANGEROUS in production setting -- pprof running", "listen-addr", PprofListenAddr)
		go func() {
			runLog.Error(http.ListenAndServe(PprofListenAddr, nil), "serving pprof")
		}()
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		return fmt.Errorf("Running manager: %s", err)
	}

	runLog.Info("Exiting")
	if server != nil {
		server.Stop()
	}

	return nil
}

// parseTLSCipherSuites tries to validate and return the user-input ciphers or returns a default list
// implementation largely stolen from: https://github.com/antrea-io/antrea/blob/25ff93d8987c6b9e3a2062254da6d7d70c623410/pkg/util/cipher/cipher.go#L32
func parseTLSCipherSuites(opts string) ([]string, error) {
	csStrList := strings.Split(strings.ReplaceAll(opts, " ", ""), ",")
	if len(csStrList) == 1 && csStrList[0] == "" {
		return nil, nil
	}

	// check to make sure they all parse - this just a fail-fast
	_, err := flag.TLSCipherSuites(csStrList)
	if err != nil {
		return nil, fmt.Errorf("unable to parse TLSCipherSuites: %s", err)
	}

	return csStrList, nil
}
