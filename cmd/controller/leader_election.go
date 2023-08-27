// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
	"os"
	"time"
)

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

func runLeaderElection(ctx context.Context, lock *resourcelock.LeaseLock, podname string, ctrlOpts Options, log logr.Logger) {
	// Start the leader election for running kapp-controller
	log.Info("Waiting for leader election")
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(c context.Context) {
				err := Run(ctrlOpts, log.WithName("controller"))
				if err != nil {
					klog.Errorf("Error while running as leader: %v", err)
				}
			},
			OnStoppedLeading: func() {
				klog.Fatalf("no longer the leader, staying inactive.")
				os.Exit(0)
			},
			OnNewLeader: func(identity string) {
				//Notify when a new leader is elected
				if identity == podname {
					return
				}
				klog.InfoS("new leader elected", "id", identity)
			},
		},
	})
}
