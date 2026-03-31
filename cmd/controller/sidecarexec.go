// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"carvel.dev/kapp-controller/pkg/sidecarexec"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func sidecarexecMain() {
	mainLog := zap.New(zap.UseDevMode(false)).WithName("kc-sidecarexec")
	mainLog.Info("start sidecarexec", "version", Version)

	reaper := sidecarexec.NewReaper(mainLog)
	go reaper.Run()

	opts := sidecarexec.ServerOpts{
		AllowedCmdNames: []string{
			// Fetch (calls impgkg and others internally)
			"vendir",
			// Template
			"ytt", "kbld", "sops", "helm", "cue",
		},
	}

	server := sidecarexec.NewServer(reaper, opts, mainLog)

	err := server.Serve()
	if err != nil {
		mainLog.Error(err, "Serving RPC")
	}
}
