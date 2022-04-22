// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os/exec"
)

// CmdRunner allows to run commands on the OS. All commands
// running within kapp-controller should happen through an instance
// of this interface so that they can be intercepted and potentially
// modified in kctrl when running kapp-controller locally.
type CmdRunner interface {
	Run(*exec.Cmd) error
	RunWithCancel(cmd *exec.Cmd, cancelCh chan struct{}) error
}

// PlainCmdRunner implements CmdRunner interface by simply running exec.Cmd.
type PlainCmdRunner struct{}

var _ CmdRunner = PlainCmdRunner{}

// NewPlainCmdRunner returns a new PlainCmdRunner.
func NewPlainCmdRunner() PlainCmdRunner {
	return PlainCmdRunner{}
}

// Run executes exec.Cmd.
func (PlainCmdRunner) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}

// RunWithCancel executes exec.Cmd.
// Kills execution immediately if value is read from cancelCh.
func (PlainCmdRunner) RunWithCancel(cmd *exec.Cmd, cancelCh chan struct{}) error {
	select {
	case <-cancelCh:
		return fmt.Errorf("Already canceled")
	default:
		// continue with execution
	}

	err := cmd.Start()
	if err != nil {
		return err
	}

	doneCh := make(chan error, 1)
	go func() {
		doneCh <- cmd.Wait()
	}()

	select {
	case <-cancelCh:
		err := cmd.Process.Kill()
		if err != nil {
			return fmt.Errorf("Killing process: %s", err)
		}
		return fmt.Errorf("Process was canceled")

	case err := <-doneCh:
		return err
	}
}
