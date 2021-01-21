// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controllerinit

// Based on https://github.com/pablo-ruth/go-init/blob/master/main.go
import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/config/global"
)

const (
	InternalControllerFlag = "internal-controller"
)

func Run(cmdName string, args []string, runLog logr.Logger) {
	runLog.Info("start init")

	go reapZombies(runLog)

	if err := configureSystem(); err != nil {
		runLog.Error(err, "Could not configure system")
		os.Exit(1)
	}

	err := runControllerCmd(cmdName, args)

	if err != nil {
		runLog.Error(err, "Could not start controller")
		os.Exit(1)
	}

	os.Exit(0)
}

func reapZombies(runLog logr.Logger) {
	runLog.Info("starting zombie reaper")

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

func runControllerCmd(cmdName string, args []string) error {
	sigs := make(chan os.Signal, 1)
	defer close(sigs)

	signal.Notify(sigs)
	defer signal.Reset()

	cmd := exec.Command(cmdName, append([]string{"--" + InternalControllerFlag}, args...)...)

	// Forward signals to child's proc group
	go func() {
		for sig := range sigs {
			if sig != syscall.SIGCHLD {
				syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal))
			}
		}
	}()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return cmd.Run()
}

func configureSystem() error {
	globalConfigurer, err := global.NewConfigurer()
	if err != nil {
		return fmt.Errorf("Creating configurer: %s", err)
	}

	err = globalConfigurer.Configure()
	if err != nil {
		return fmt.Errorf("Applying configuration: %s", err)
	}

	return nil
}
