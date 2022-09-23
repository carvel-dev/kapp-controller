// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"fmt"
	"net/http"         // Pprof related
	_ "net/http/pprof" // Pprof related
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/cmd/controller/handlers"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // Initialize gcp client auth plugin
	"k8s.io/component-base/cli/flag"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"

	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	kcconfig "github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"

	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	pkgclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned"
)

const (
	PprofListenAddr       = "0.0.0.0:6060"
	kappctrlAPIPORTEnvKey = "KAPPCTRL_API_PORT"
)

type Options struct {
	Concurrency       int
	Namespace         string
	EnablePprof       bool
	APIRequestTimeout time.Duration
	PackagingGloablNS string
	TLSCipherSuites   string
}

// Based on https://github.com/kubernetes-sigs/controller-runtime/blob/8f633b179e1c704a6e40440b528252f147a3362a/examples/builtins/main.go
func Run(opts Options, runLog logr.Logger) {
	runLog.Info("start controller")
	runLog.Info("setting up manager")

	restConfig := config.GetConfigOrDie()

	if opts.APIRequestTimeout != 0 {
		restConfig.Timeout = opts.APIRequestTimeout
	}

	mgr, err := manager.New(restConfig, manager.Options{Namespace: opts.Namespace, Scheme: kcconfig.Scheme})
	if err != nil {
		runLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	logProxies(runLog)

	runLog.Info("setting up controller")

	coreClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		runLog.Error(err, "building core client")
		os.Exit(1)
	}

	kcClient, err := kcclient.NewForConfig(restConfig)
	if err != nil {
		runLog.Error(err, "building app client")
		os.Exit(1)
	}

	kcConfig, err := kcconfig.GetConfig(coreClient)
	if err != nil {
		runLog.Error(err, "getting kapp-controller config")
		os.Exit(1)
	}

	pkgClient, err := pkgclient.NewForConfig(restConfig)
	if err != nil {
		runLog.Error(err, "building app client")
		os.Exit(1)
	}

	// assign bindPort to env var KAPPCTRL_API_PORT if available
	var bindPort int
	if apiPort, ok := os.LookupEnv(kappctrlAPIPORTEnvKey); ok {
		var err error
		if bindPort, err = strconv.Atoi(apiPort); err != nil {
			runLog.Error(fmt.Errorf("%s environment variable must be an integer", kappctrlAPIPORTEnvKey), "reading server port")
			os.Exit(1)
		}
	} else {
		runLog.Error(fmt.Errorf("os call failed to read env var %s", kappctrlAPIPORTEnvKey), "reading server port")
		os.Exit(1)
	}

	cSuites, err := parseTLSCipherSuites(opts.TLSCipherSuites)
	if err != nil {
		runLog.Error(err, "creating API server %s", err)
		os.Exit(1)
	}

	server, err := apiserver.NewAPIServer(restConfig, coreClient, kcClient, opts.PackagingGloablNS, bindPort, cSuites)
	if err != nil {
		runLog.Error(err, "creating API server %s", err)
		os.Exit(1)
	}
	err = server.Run()
	if err != nil {
		runLog.Error(err, "starting server")
		os.Exit(1)
	}

	refTracker := reftracker.NewAppRefTracker()
	updateStatusTracker := reftracker.NewAppUpdateStatus()

	appFactory := AppFactory{
		coreClient: coreClient,
		kcConfig:   kcConfig,
		appClient:  kcClient,
	}

	{ // add controller for apps
		schApp := handlers.NewSecretHandler(runLog, refTracker, updateStatusTracker)
		cfgmhApp := handlers.NewConfigMapHandler(runLog, refTracker, updateStatusTracker)
		ctrlAppOpts := controller.Options{
			Reconciler: NewUniqueReconciler(&ErrReconciler{
				delegate: NewAppsReconciler(kcClient, runLog.WithName("ar"), appFactory, refTracker, updateStatusTracker),
				log:      runLog.WithName("pr"),
			}),
			MaxConcurrentReconciles: opts.Concurrency,
		}

		ctrlApp, err := controller.New("kapp-controller-app", mgr, ctrlAppOpts)
		if err != nil {
			runLog.Error(err, "unable to set up kapp-controller-app")
			os.Exit(1)
		}

		err = ctrlApp.Watch(&source.Kind{Type: &kcv1alpha1.App{}}, &handler.EnqueueRequestForObject{})
		if err != nil {
			runLog.Error(err, "unable to watch Apps")
			os.Exit(1)
		}

		err = ctrlApp.Watch(&source.Kind{Type: &v1.Secret{}}, schApp)
		if err != nil {
			runLog.Error(err, "unable to watch Secrets")
			os.Exit(1)
		}

		err = ctrlApp.Watch(&source.Kind{Type: &v1.ConfigMap{}}, cfgmhApp)
		if err != nil {
			runLog.Error(err, "unable to watch ConfigMaps")
			os.Exit(1)
		}
	}

	{ // add controller for PackageInstall
		pkgInstallCtrlOpts := controller.Options{
			Reconciler: &PackageInstallReconciler{
				kcClient:  kcClient,
				pkgClient: pkgClient,
				log:       runLog.WithName("ipr"),
			},
			MaxConcurrentReconciles: opts.Concurrency,
		}

		pkgInstallCtrl, err := controller.New("kapp-controller-packageinstall", mgr, pkgInstallCtrlOpts)
		if err != nil {
			runLog.Error(err, "unable to set up kapp-controller-packageinstall")
			os.Exit(1)
		}

		err = pkgInstallCtrl.Watch(&source.Kind{Type: &pkgingv1alpha1.PackageInstall{}}, &handler.EnqueueRequestForObject{})
		if err != nil {
			runLog.Error(err, "unable to watch *pkgingv1alpha1.PackageInstall")
			os.Exit(1)
		}

		err = pkgInstallCtrl.Watch(&source.Kind{Type: &datapkgingv1alpha1.Package{}}, handlers.NewPackageInstallVersionHandler(kcClient, opts.PackagingGloablNS, runLog.WithName("handler")))
		if err != nil {
			runLog.Error(err, "unable to watch *datapkgingv1alpha1.Package for PackageInstall")
			os.Exit(1)
		}

		err = pkgInstallCtrl.Watch(&source.Kind{Type: &kcv1alpha1.App{}}, &handler.EnqueueRequestForOwner{
			OwnerType:    &pkgingv1alpha1.PackageInstall{},
			IsController: true,
		})
		if err != nil {
			runLog.Error(err, "unable to watch *kcv1alpha1.App for PackageInstall")
			os.Exit(1)
		}
	}

	{ // add controller for pkgrepositories
		schRepo := handlers.NewSecretHandler(runLog, refTracker, updateStatusTracker)

		pkgRepositoriesCtrlOpts := controller.Options{
			Reconciler: NewPkgRepositoryReconciler(kcClient, runLog.WithName("prr"), appFactory, refTracker, updateStatusTracker),
			// TODO: Consider making this configurable for multiple PackageRepo reconciles
			MaxConcurrentReconciles: 1,
		}

		pkgRepositoryCtrl, err := controller.New("kapp-controller-package-repository", mgr, pkgRepositoriesCtrlOpts)
		if err != nil {
			runLog.Error(err, "unable to set up kapp-controller-package-repository")
			os.Exit(1)
		}

		err = pkgRepositoryCtrl.Watch(&source.Kind{Type: &pkgingv1alpha1.PackageRepository{}}, &handler.EnqueueRequestForObject{})
		if err != nil {
			runLog.Error(err, "unable to watch *pkgingv1alpha1.PackageRepository")
			os.Exit(1)
		}

		err = pkgRepositoryCtrl.Watch(&source.Kind{Type: &v1.Secret{}}, schRepo)
		if err != nil {
			runLog.Error(err, "unable to watch Secrets")
			os.Exit(1)
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
		runLog.Error(err, "unable to run manager")
		os.Exit(1)
	}

	runLog.Info("Exiting")
	server.Stop()
	os.Exit(0)
}

func logProxies(runLog logr.Logger) {
	if proxyVal := os.Getenv("http_proxy"); proxyVal != "" {
		runLog.Info(fmt.Sprintf("Using http proxy '%s'", proxyVal))
	}

	if proxyVal := os.Getenv("https_proxy"); proxyVal != "" {
		runLog.Info(fmt.Sprintf("Using https proxy '%s'", proxyVal))
	}

	if noProxyVal := os.Getenv("no_proxy"); noProxyVal != "" {
		runLog.Info(fmt.Sprintf("No proxy set for: %s", noProxyVal))
	}
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
