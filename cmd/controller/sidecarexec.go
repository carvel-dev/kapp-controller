// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	goexec "os/exec"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/sidecarexec"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func sidecarexecMain(debug bool, debugArgs []string) {
	localCmdRunner := exec.NewPlainCmdRunner()
	sandboxCmdRunner := sidecarexec.NewSandboxCmdRunner(localCmdRunner, sidecarexec.SandboxCmdRunnerOpts{
		RequiresPosix: map[string]bool{
			"vendir": true,
			"bash":   true, // for debugging
		},
		RequiresNetwork: map[string]bool{
			"vendir": true,
			"kbld":   true,
		},
	})

	if debug {
		cmd := goexec.Command("bash")
		if len(debugArgs) > 0 {
			cmd = goexec.Command(debugArgs[0], debugArgs[1:]...)
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := sandboxCmdRunner.Run(cmd, exec.RunOpts{})
		if err != nil {
			fmt.Printf("Exit error: %s\n", err)
		}
		return
	}

	mainLog := zap.New(zap.UseDevMode(false)).WithName("kc-sidecarexec")
	mainLog.Info("start sidecarexec", "version", Version)

	go reapZombies(mainLog)

	serverOpts := sidecarexec.ServerOpts{
		AllowedCmdNames: []string{
			// Fetch (calls impgkg and others internally)
			"vendir",
			// Template
			"ytt", "kbld", "sops", "helm", "cue",
		},
	}
	server := sidecarexec.NewServer(sandboxCmdRunner, serverOpts, mainLog)

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
