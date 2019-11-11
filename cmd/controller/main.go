package main

// Based on https://github.com/kubernetes-sigs/controller-runtime/blob/8f633b179e1c704a6e40440b528252f147a3362a/examples/builtins/main.go

import (
	"flag"
	"os"

	kcv1alpha1 "github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/k14s/kapp-controller/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	log             = logf.Log.WithName("kapp-controller")
	ctrlConcurrency = 10
	ctrlNamespace   = ""
)

func main() {
	flag.IntVar(&ctrlConcurrency, "concurrency", 10, "Max concurrent reconciles")
	flag.StringVar(&ctrlNamespace, "namespace", "", "Namespace to watch")
	flag.Parse()

	logf.SetLogger(zap.Logger(false))
	entryLog := log.WithName("entrypoint")

	entryLog.Info("setting up manager")

	restConfig := config.GetConfigOrDie()

	mgr, err := manager.New(restConfig, manager.Options{Namespace: ctrlNamespace})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	entryLog.Info("Setting up controllers")

	coreClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		entryLog.Error(err, "building core client")
		os.Exit(1)
	}

	appClient, err := kcclient.NewForConfig(restConfig)
	if err != nil {
		entryLog.Error(err, "building app client")
		os.Exit(1)
	}

	appFactory := AppFactory{
		coreClient: coreClient,
		appClient:  appClient,
	}

	{ // add controller for apps
		ctrlAppOpts := controller.Options{
			Reconciler: &PeriodicReconciler{NewUniqueReconciler(&AppsReconciler{
				client:     mgr.GetClient(),
				appFactory: appFactory,
				log:        log.WithName("ar"),
			})},
			MaxConcurrentReconciles: ctrlConcurrency,
		}

		ctrlApp, err := controller.New("kapp-controller-app", mgr, ctrlAppOpts)
		if err != nil {
			entryLog.Error(err, "unable to set up kapp-controller-app")
			os.Exit(1)
		}

		err = ctrlApp.Watch(&source.Kind{Type: &kcv1alpha1.App{}}, &handler.EnqueueRequestForObject{})
		if err != nil {
			entryLog.Error(err, "unable to watch *kcv1alpha1.App")
			os.Exit(1)
		}
	}

	entryLog.Info("starting manager")

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
