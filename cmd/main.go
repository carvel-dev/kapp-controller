// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"os"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/cmd/controller"
	"github.com/vmware-tanzu/carvel-kapp-controller/cmd/controllerinit"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Version of kapp-controller is set via ldflags at build-time from the most recent git tag; see hack/build.sh
var Version = "develop"

func main() {
	ctrlOpts := controller.Options{}

	var runController bool

	flag.IntVar(&ctrlOpts.Concurrency, "concurrency", 10, "Max concurrent reconciles")
	flag.StringVar(&ctrlOpts.Namespace, "namespace", "", "Namespace to watch")
	flag.StringVar(&ctrlOpts.PackagingGloablNS, "packaging-global-namespace", "", "The namespace used for global packaging resources")
	flag.StringVar(&ctrlOpts.MetricsBindAddress, "metrics-bind-address", ":8080", "Address for metrics server. If 0, then metrics server doesnt listen on any port.")
	flag.BoolVar(&ctrlOpts.EnablePprof, "dangerous-enable-pprof", false, "If set to true, enable pprof on "+controller.PprofListenAddr)
	flag.DurationVar(&ctrlOpts.APIRequestTimeout, "api-request-timeout", time.Duration(0), "HTTP timeout for Kubernetes API requests")
	flag.BoolVar(&runController, controllerinit.InternalControllerFlag, false, "[Internal] run the controller code")
	flag.BoolVar(&ctrlOpts.APIPriorityAndFairness, "enable-APIPriorityAndFairness", true, "Enable/disable APIPriorityAndFairness feature gate for apiserver. Recommended to disable for <= k8s 1.19.")
	flag.Parse()

	log := logf.Log.WithName("kc")

	logf.SetLogger(zap.New(zap.UseDevMode(false)))

	mainLog := log.WithName("main")
	mainLog.Info("kapp-controller", "version", Version)

	if runController {
		err := controller.Run(ctrlOpts, log.WithName("controller"))
		if err != nil {
			mainLog.Error(err, "Exited run with error")
			os.Exit(1)
		}
		os.Exit(0)
		panic("unreachable: controller returned")
	}

	controllerinit.Run(os.Args[0], os.Args[1:], log.WithName("init"))
	panic("unreachable: init proc returned")
}
