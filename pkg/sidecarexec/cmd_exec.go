// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"bytes"
	"fmt"
	"os"
	goexec "os/exec"

	"carvel.dev/kapp-controller/pkg/exec"
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
	local           exec.CmdRunner
	allowedCmdNames map[string]struct{}
}

// Run executes a command (out of a set of allowed ones).
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

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := r.local.Run(cmd)
	if err != nil {
		output.Error = err.Error()
		output.ExitCode = -1
		if exitError, ok := err.(*goexec.ExitError); ok {
			output.ExitCode = exitError.ExitCode()
		}
	}

	output.Stdout = stdout.Bytes()
	output.Stderr = stderr.Bytes()
	return nil
}
