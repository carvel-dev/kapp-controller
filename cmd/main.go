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

const (
	Version = "0.17.0-alpha.1"
)

func main() {
	ctrlOpts := controller.Options{}

	var runController bool

	flag.IntVar(&ctrlOpts.Concurrency, "concurrency", 10, "Max concurrent reconciles")
	flag.StringVar(&ctrlOpts.Namespace, "namespace", "", "Namespace to watch")
	flag.BoolVar(&ctrlOpts.EnablePprof, "dangerous-enable-pprof", false, "If set to true, enable pprof on "+controller.PprofListenAddr)
	flag.DurationVar(&ctrlOpts.APIRequestTimeout, "api-request-timeout", time.Duration(0), "HTTP timeout for Kubernetes API requests")
	flag.BoolVar(&runController, controllerinit.InternalControllerFlag, false, "[Internal] run the controller code")
	flag.Parse()

	log := logf.Log.WithName("kc")

	logf.SetLogger(zap.Logger(false))

	mainLog := log.WithName("main")
	mainLog.Info("kapp-controller", "version", Version)

	if runController {
		controller.Run(ctrlOpts, log.WithName("controller"))
		panic("unreachable: controller returned")
	}

	controllerinit.Run(os.Args[0], os.Args[1:], log.WithName("init"))
	panic("unreachable: init proc returned")
}
