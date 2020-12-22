// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package initproc

// Based on https://github.com/pablo-ruth/go-init/blob/master/main.go
import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-logr/logr"
)

func Run(cmdName string, args []string, runLog logr.Logger) {
	runLog.Info("start init")

	zombiesCtx, zombiesCancelFunc := context.WithCancel(context.Background())

	var zombiesWg sync.WaitGroup
	zombiesWg.Add(1)

	go reapZombies(zombiesCtx, &zombiesWg, runLog)

	err := runControllerCmd(cmdName, args)

	zombiesCancelFunc()
	zombiesWg.Wait()

	if err != nil {
		runLog.Error(err, "Could not start controller")
		os.Exit(1)
	}

	os.Exit(0)
}

func reapZombies(ctx context.Context, zombiesWg *sync.WaitGroup, runLog logr.Logger) {
	runLog.Info("starting zombie reaper")

	for {
		var status syscall.WaitStatus

		pid, _ := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)
		if pid <= 0 {
			time.Sleep(1 * time.Second)
		} else {
			continue
		}

		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
		}
	}
}

func runControllerCmd(cmdName string, args []string) (int, error) {
	sigs := make(chan os.Signal, 1)
	defer close(sigs)

	signal.Notify(sigs)
	defer signal.Reset()

	// Forward signals to child's proc group
	go func() {
		for sig := range sigs {
			if sig != syscall.SIGCHLD {
				syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal))
			}
		}
	}()

	cmd := exec.Command(cmdName, append([]string{"--internal-controller"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err := cmd.Run()
	if err != nil {
		return 1, err
	}

	return 0, nil
}
