// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"fmt"
	"io/ioutil"
	goexec "os/exec"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

// CmdExecClient executes commands remotely (in a sidecar container)
// except for kapp commands which continue to run locally.
type CmdExecClient struct {
	local     exec.CmdRunner
	rpcClient rpcClient
}

var _ exec.CmdRunner = CmdExecClient{}

type runError struct {
	exitCode int
	message  string
}

func (e runError) ExitCode() int { return e.exitCode }
func (e runError) Error() string { return e.message }

// Run makes a CmdExec.Run RPC call. kapp command run locally though.
func (r CmdExecClient) Run(cmd *goexec.Cmd) error {
	// TODO is there better way to "undo" path resolution done by exec.Command
	cmdName := filepath.Base(cmd.Path)
	args := cmd.Args[1:]

	if cmdName == "kapp" {
		return r.local.Run(cmd)
	}

	input := CmdInput{
		Command: cmdName,
		Args:    args,
		Env:     cmd.Env,
		Dir:     cmd.Dir,
	}

	if cmd.Stdin != nil {
		bs, err := ioutil.ReadAll(cmd.Stdin)
		if err != nil {
			return fmt.Errorf("Reading stdin: %s", err)
		}
		input.Stdin = bs
	}

	var output CmdOutput

	err := r.rpcClient.Call("CmdExec.Run", input, &output)
	if err != nil {
		return fmt.Errorf("Internal run comm: %s", err)
	}

	if cmd.Stdout != nil {
		cmd.Stdout.Write(output.Stdout)
	}
	if cmd.Stderr != nil {
		cmd.Stderr.Write(output.Stderr)
	}

	if output.ExitCode != 0 || len(output.Error) > 0 {
		return runError{exitCode: output.ExitCode, message: output.Error}
	}
	return nil
}

// RunWithCancel is not supported except for kapp which runs locally.
func (r CmdExecClient) RunWithCancel(cmd *goexec.Cmd, cancelCh chan struct{}) error {
	cmdName := filepath.Base(cmd.Path)

	if cmdName == "kapp" {
		return r.local.RunWithCancel(cmd, cancelCh)
	}

	panic("Internal inconsistency: RunWithCancel not supported")
}
