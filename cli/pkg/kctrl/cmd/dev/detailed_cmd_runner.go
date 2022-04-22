// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package dev

import (
	"fmt"
	"io"
	goexec "os/exec"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

type DetailedCmdRunner struct {
	log        io.Writer
	fullOutput bool
}

var _ exec.CmdRunner = &DetailedCmdRunner{}

func NewDetailedCmdRunner(log io.Writer, fullOutput bool) *DetailedCmdRunner {
	return &DetailedCmdRunner{log, fullOutput}
}

func (r DetailedCmdRunner) Run(cmd *goexec.Cmd) error {
	if r.fullOutput {
		cmd.Stdout = io.MultiWriter(r.log, cmd.Stdout)
		cmd.Stderr = io.MultiWriter(r.log, cmd.Stderr)
	}

	fmt.Fprintf(r.log, "==> Executing %s %v\n", cmd.Path, cmd.Args)
	defer fmt.Fprintf(r.log, "==> Finished executing %s\n\n", cmd.Path)

	return exec.PlainCmdRunner{}.Run(cmd)
}

func (r DetailedCmdRunner) RunWithCancel(cmd *goexec.Cmd, cancelCh chan struct{}) error {
	if r.fullOutput {
		cmd.Stdout = io.MultiWriter(r.log, cmd.Stdout)
		cmd.Stderr = io.MultiWriter(r.log, cmd.Stderr)
	}

	fmt.Fprintf(r.log, "==> Executing %s %v\n", cmd.Path, cmd.Args)
	defer fmt.Fprintf(r.log, "==> Finished executing %s\n\n", cmd.Path)

	return exec.PlainCmdRunner{}.RunWithCancel(cmd, cancelCh)
}
