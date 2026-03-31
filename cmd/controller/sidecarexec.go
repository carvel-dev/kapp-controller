// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"carvel.dev/kapp-controller/pkg/exec"
	"carvel.dev/kapp-controller/pkg/sidecarexec"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func sidecarexecMain() {
	mainLog := zap.New(zap.UseDevMode(false)).WithName("kc-sidecarexec")
	mainLog.Info("start sidecarexec", "version", Version)

	// Note: Zombie reaping is now handled by tini as PID 1

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
