// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"syscall"
	"time"

	"carvel.dev/kapp-controller/pkg/exec"
	"carvel.dev/kapp-controller/pkg/sidecarexec"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func sidecarexecMain() {
	mainLog := zap.New(zap.UseDevMode(false)).WithName("kc-sidecarexec")
	mainLog.Info("start sidecarexec", "version", Version)

	go reapZombies(mainLog)

	localCmdRunner := exec.NewPlainCmdRunner()
	opts := sidecarexec.ServerOpts{
		AllowedCmdNames: []string{
			// Fetch (calls impgkg and others internally)
			"vendir",
			// Template
			"ytt", "kbld", "sops", "helm", "cue",
		},
	}

	server := sidecarexec.NewServer(localCmdRunner, opts, mainLog)

	err := server.Serve()
	if err != nil {
		mainLog.Error(err, "Serving RPC")
	}
}

func reapZombies(log logr.Logger) {
	log.Info("starting zombie reaper")

	for {
		var status syscall.WaitStatus

		pid, _ := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)
		if pid <= 0 {
			time.Sleep(1 * time.Second)
		} else {
			continue
		}
	}
}
