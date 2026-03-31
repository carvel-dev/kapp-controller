// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	goexec "os/exec"
	"sync"
)

// CmdInput describes a command to run.
type CmdInput struct {
	Command string
	Args    []string
	Stdin   []byte
	Env     []string
	Dir     string
}

// CmdOutput describes an command execution result.
type CmdOutput struct {
	Stdout   []byte
	Stderr   []byte
	Error    string
	ExitCode int
}

// CmdExec provides RPC interface for command execution.
type CmdExec struct {
	reaper          *Reaper
	allowedCmdNames map[string]struct{}
}

// Run executes a command (out of a set of allowed ones).
// It uses cmd.Start() and delegates process wait to the centralized Reaper
// to avoid a race with the zombie reaping loop.
func (r CmdExec) Run(input CmdInput, output *CmdOutput) error {
	if _, found := r.allowedCmdNames[input.Command]; !found {
		return fmt.Errorf("Command '%s' is not allowed", input.Command)
	}

	cmd := goexec.Command(input.Command, input.Args...)

	if len(input.Stdin) > 0 {
		cmd.Stdin = bytes.NewBuffer(input.Stdin)
	}
	cmd.Env = os.Environ()
	if len(input.Env) > 0 {
		cmd.Env = append(cmd.Env, input.Env...)
	}
	if len(input.Dir) > 0 {
		cmd.Dir = input.Dir
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Creating stdout pipe: %s", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Creating stderr pipe: %s", err)
	}

	// Hold the reaper lock across Start + Track to guarantee the reaping
	// loop cannot consume this child's exit status before we register it.
	r.reaper.Lock()
	if err := cmd.Start(); err != nil {
		r.reaper.Unlock()
		return fmt.Errorf("Starting command: %s", err)
	}
	statusCh := r.reaper.TrackLocked(cmd.Process.Pid)
	r.reaper.Unlock()

	var stdout, stderr bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(&stdout, stdoutPipe)
	}()
	go func() {
		defer wg.Done()
		io.Copy(&stderr, stderrPipe)
	}()

	// Wait for all output to be read before checking exit status.
	// The pipes close when the process exits, so this will unblock.
	wg.Wait()

	// Now wait for the Reaper to deliver the exit status.
	waitStatus := <-statusCh

	// Tell Go's os/exec we handled the wait ourselves.
	cmd.Process.Release()

	output.Stdout = stdout.Bytes()
	output.Stderr = stderr.Bytes()

	if waitStatus.Signaled() {
		output.ExitCode = -1
		output.Error = fmt.Sprintf("process killed by signal %d", waitStatus.Signal())
	} else if waitStatus.Exited() && waitStatus.ExitStatus() != 0 {
		output.ExitCode = waitStatus.ExitStatus()
		output.Error = fmt.Sprintf("exit status %d", waitStatus.ExitStatus())
	}

	return nil
}
