// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"fmt"
	"io/ioutil"
	"net/rpc"
	goexec "os/exec"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

// Client executes commands remotely (in a sidecar container)
// except for kapp commands which continue to run locally.
type Client struct {
	local     exec.CmdRunner
	rpcClient *rpc.Client
}

var _ exec.CmdRunner = Client{}

// NewClient returns a new Client.
func NewClient(local exec.CmdRunner) Client {
	rpcClient, err := rpc.DialHTTP(serverListenType, serverListenAddr)
	if err != nil {
		// TODO
		panic("Dialing rpc")
	}

	return Client{local, rpcClient}
}

func (r Client) Run(cmd *goexec.Cmd) error {
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

	err := r.rpcClient.Call("ServerMethods.Run", input, &output)
	if err != nil {
		return fmt.Errorf("Internal run comm: %s", err)
	}

	if cmd.Stdout != nil {
		cmd.Stdout.Write(output.Stdout)
	}
	if cmd.Stderr != nil {
		cmd.Stderr.Write(output.Stderr)
	}

	// TODO exit code on the error
	if len(output.Error) > 0 {
		return fmt.Errorf("%s", output.Error)
	}
	return nil
}

func (r Client) RunWithCancel(cmd *goexec.Cmd, cancelCh chan struct{}) error {
	cmdName := filepath.Base(cmd.Path)

	if cmdName == "kapp" {
		return r.local.RunWithCancel(cmd, cancelCh)
	}

	panic("Internal inconsistency: RunWithCancel not supported")
}
