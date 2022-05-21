// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"bytes"
	"fmt"
	goexec "os/exec"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

type CmdInput struct {
	Command string
	Args    []string
	Stdin   []byte
	Env     []string
	Dir     string
}

type CmdOutput struct {
	Stdout []byte
	Stderr []byte
	Error  string
}

type CmdExec struct {
	local           exec.CmdRunner
	allowedCmdNames map[string]struct{}
}

func (r CmdExec) Run(input CmdInput, output *CmdOutput) error {
	if _, found := r.allowedCmdNames[input.Command]; !found {
		return fmt.Errorf("Command '%s' is not allowed", input.Command)
	}

	cmd := goexec.Command(input.Command, input.Args...)

	if len(input.Stdin) > 0 {
		cmd.Stdin = bytes.NewBuffer(input.Stdin)
	}
	if len(input.Env) > 0 {
		cmd.Env = input.Env
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
	}

	output.Stdout = stdout.Bytes()
	output.Stderr = stderr.Bytes()
	return nil
}
