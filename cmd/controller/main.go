// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/sidecarexec"
	"k8s.io/klog/v2"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Version of kapp-controller is set via ldflags at build-time from the most recent git tag; see hack/build.sh
var Version = "develop"

func main() {
	ctrlOpts := Options{}
	var isSidecarexec, isSidecarexecDebug, isSidecarexecWrap bool

	flag.IntVar(&ctrlOpts.Concurrency, "concurrency", 10, "Max concurrent reconciles")
	flag.StringVar(&ctrlOpts.Namespace, "namespace", "", "Namespace to watch")
	flag.StringVar(&ctrlOpts.PackagingGloablNS, "packaging-global-namespace", "", "The namespace used for global packaging resources")
	flag.StringVar(&ctrlOpts.MetricsBindAddress, "metrics-bind-address", ":8080", "Address for metrics server. If 0, then metrics server doesnt listen on any port.")
	flag.BoolVar(&ctrlOpts.EnablePprof, "dangerous-enable-pprof", false, "If set to true, enable pprof on "+PprofListenAddr)
	flag.DurationVar(&ctrlOpts.APIRequestTimeout, "api-request-timeout", time.Duration(0), "HTTP timeout for Kubernetes API requests")
	flag.BoolVar(&ctrlOpts.APIPriorityAndFairness, "enable-api-priority-and-fairness", true, "Enable/disable APIPriorityAndFairness feature gate for apiserver. Recommended to disable for <= k8s 1.19.")
	flag.BoolVar(&isSidecarexec, "sidecarexec", false, "Run sidecarexec")
	flag.BoolVar(&isSidecarexecDebug, "sidecarexecdebug", false, "Run sidecarexecdebug")
	flag.BoolVar(&isSidecarexecWrap, "sidecarexecwrap", false, "Run sidecarexecwrap")
	flag.Parse()

	if isSidecarexec || isSidecarexecDebug {
		sidecarexecMain(isSidecarexecDebug, flag.Args())
		return
	}
	if isSidecarexecWrap {
		err := sidecarexec.SandboxWrap{}.ExecuteCmd(flag.Args())
		if err != nil {
			fmt.Fprintf(os.Stderr, "sidecarexecwrap: Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
		return
	}

	log := zap.New(zap.UseDevMode(false)).WithName("kc")
	logf.SetLogger(log)
	klog.SetLogger(log)

	mainLog := log.WithName("main")
	mainLog.Info("kapp-controller", "version", Version)

	err := Run(ctrlOpts, log.WithName("controller"))
	if err != nil {
		mainLog.Error(err, "Exited run with error")
		os.Exit(1)
	}

	os.Exit(0)
}
