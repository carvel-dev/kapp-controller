// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

// Based on https://github.com/kubernetes-sigs/controller-runtime/blob/8f633b179e1c704a6e40440b528252f147a3362a/examples/builtins/main.go

import (
	"os"

	// Pprof related
	"net/http"
	_ "net/http/pprof"

	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func RunController(ctrlConcurrency int, ctrlNamespace string, enablePprof bool, pprofListenAddr string, controllerLog logr.Logger) {
	controllerLog.Info("setting up manager")

	restConfig := config.GetConfigOrDie()

	mgr, err := manager.New(restConfig, manager.Options{Namespace: ctrlNamespace})
	if err != nil {
		controllerLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	controllerLog.Info("setting up controller")

	coreClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		controllerLog.Error(err, "building core client")
		os.Exit(1)
	}

	appClient, err := kcclient.NewForConfig(restConfig)
	if err != nil {
		controllerLog.Error(err, "building app client")
		os.Exit(1)
	}

	appFactory := AppFactory{
		coreClient: coreClient,
		appClient:  appClient,
	}

	{ // add controller for apps
		ctrlAppOpts := controller.Options{
			Reconciler: NewUniqueReconciler(&ErrReconciler{
				delegate: &AppsReconciler{
					appClient:  appClient,
					appFactory: appFactory,
					log:        controllerLog.WithName("ar"),
				},
				log: controllerLog.WithName("pr"),
			}),
			MaxConcurrentReconciles: ctrlConcurrency,
		}

		ctrlApp, err := controller.New("kapp-controller-app", mgr, ctrlAppOpts)
		if err != nil {
			controllerLog.Error(err, "unable to set up kapp-controller-app")
			os.Exit(1)
		}

		err = ctrlApp.Watch(&source.Kind{Type: &kcv1alpha1.App{}}, &handler.EnqueueRequestForObject{})
		if err != nil {
			controllerLog.Error(err, "unable to watch *kcv1alpha1.App")
			os.Exit(1)
		}
	}

	controllerLog.Info("starting manager")

	if enablePprof {
		controllerLog.Info("DANGEROUS in production setting -- pprof running", "listen-addr", pprofListenAddr)
		go func() {
			controllerLog.Error(http.ListenAndServe(pprofListenAddr, nil), "serving pprof")
		}()
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		controllerLog.Error(err, "unable to run manager")
		os.Exit(1)
	}

	controllerLog.Info("Exiting")
	os.Exit(0)
}
