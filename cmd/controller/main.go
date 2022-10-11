// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Version of kapp-controller is set via ldflags at build-time from the most recent git tag; see hack/build.sh
var Version = "develop"

var (
	client *clientset.Clientset
)

func getNewLock(lockname, podname, namespace string) *resourcelock.LeaseLock {
	return &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      lockname,
			Namespace: namespace,
		},
		Client: client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: podname,
		},
	}
}

func runLeaderElection(lock *resourcelock.LeaseLock, ctx context.Context, id string, ctrlOpts Options) {
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(c context.Context) {
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
			},
			OnStoppedLeading: func() {
				klog.Info("no longer the leader, staying inactive.")
			},
			OnNewLeader: func(current_id string) {
				if current_id == id {
					klog.Info("still the leader!")
					return
				}
				klog.Info("new leader is %s", current_id)
			},
		},
	})
}

func main() {
	ctrlOpts := Options{}
	var sidecarexec bool

	flag.IntVar(&ctrlOpts.Concurrency, "concurrency", 10, "Max concurrent reconciles")
	flag.StringVar(&ctrlOpts.Namespace, "namespace", "", "Namespace to watch")
	flag.StringVar(&ctrlOpts.PackagingGlobalNS, "packaging-global-namespace", "", "The namespace used for global packaging resources")
	flag.StringVar(&ctrlOpts.MetricsBindAddress, "metrics-bind-address", ":8080", "Address for metrics server. If 0, then metrics server doesnt listen on any port.")
	flag.BoolVar(&ctrlOpts.EnablePprof, "dangerous-enable-pprof", false, "If set to true, enable pprof on "+PprofListenAddr)
	flag.DurationVar(&ctrlOpts.APIRequestTimeout, "api-request-timeout", time.Duration(0), "HTTP timeout for Kubernetes API requests")
	flag.BoolVar(&ctrlOpts.APIPriorityAndFairness, "enable-api-priority-and-fairness", true, "Enable/disable APIPriorityAndFairness feature gate for apiserver. Recommended to disable for <= k8s 1.19.")
	flag.BoolVar(&sidecarexec, "sidecarexec", false, "Run sidecarexec")
	flag.BoolVar(&ctrlOpts.StartAPIServer, "start-api-server", true, "Start apiserver")
	flag.StringVar(&ctrlOpts.TLSCipherSuites, "tls-cipher-suites", "", "comma separated list of acceptable cipher suites. Empty list will use defaults from underlying libraries.")
	flag.Parse()

	if sidecarexec {
		sidecarexecMain()
		return
	}

	var (
		leaseLockName      string
		leaseLockNamespace string
		podName            = os.Getenv("POD_NAME")
	)
	config, err := rest.InClusterConfig()
	client = clientset.NewForConfigOrDie(config)

	if err != nil {
		klog.Fatalf("failed to get kubeconfig")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lock := getNewLock(leaseLockName, podName, leaseLockNamespace)
	runLeaderElection(lock, ctx, podName, ctrlOpts)

	os.Exit(0)
}
