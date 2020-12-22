package main

import (
	"flag"
	"os"

	controller "github.com/vmware-tanzu/carvel-kapp-controller/cmd/controller"
	"github.com/vmware-tanzu/carvel-kapp-controller/cmd/initproc"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	ctrlConcurrency = 10
	ctrlNamespace   = ""
	enablePprof     = false
	runController   = false
	log             = logf.Log.WithName("kc")
)

const (
	pprofListenAddr = "0.0.0.0:6060"
	Version         = "0.13.0"
)

// TODO: add logging
func main() {
	flag.BoolVar(&runController, "c", false, "Run the controller code")
	flag.IntVar(&ctrlConcurrency, "concurrency", 10, "Max concurrent reconciles")
	flag.StringVar(&ctrlNamespace, "namespace", "", "Namespace to watch")
	flag.BoolVar(&enablePprof, "dangerous-enable-pprof", false, "If set to true, enable pprof on "+pprofListenAddr)
	flag.Parse()

	logf.SetLogger(zap.Logger(false))
	entryLog := log.WithName("entrypoint")
	entryLog.Info("kapp-controller", "version", Version)

	if runController {
		entryLog.Info("running controller")
		controller.RunController(ctrlConcurrency, ctrlNamespace, enablePprof, pprofListenAddr, log.WithName("controller"))

		// unreachable
		panic("controller returned")
	}

	// preserve flags for exec of controller
	entryLog.Info("running init")
	initproc.RunInit(os.Args[1:], log.WithName("init"))

	// unreachable
	panic("init proc returned")

}
