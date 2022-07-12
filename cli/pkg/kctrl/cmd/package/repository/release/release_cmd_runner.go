// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"fmt"
	"io"
	goexec "os/exec"
	"strings"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

type ReleaseCmdRunner struct {
	log        io.Writer
	fullOutput bool

	tempImgLockFilepath string
}

var _ exec.CmdRunner = &ReleaseCmdRunner{}

func NewReleaseCmdRunner(log io.Writer, fullOutput bool, tempImgLockFilepath string) *ReleaseCmdRunner {
	return &ReleaseCmdRunner{log, fullOutput, tempImgLockFilepath}
}

func (r ReleaseCmdRunner) Run(cmd *goexec.Cmd) error {
	if strings.Contains(cmd.Path, "/kapp") {
		return nil
	}
	if strings.Contains(cmd.Path, "/kbld") {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--imgpkg-lock-output=%s", r.tempImgLockFilepath))
	}
	if r.fullOutput {
		cmd.Stdout = io.MultiWriter(r.log, cmd.Stdout)
		cmd.Stderr = io.MultiWriter(r.log, cmd.Stderr)
	}

	fmt.Fprintf(r.log, "==> Executing %s %v\n", cmd.Path, cmd.Args)
	defer fmt.Fprintf(r.log, "==> Finished executing %s\n\n", cmd.Path)

	return exec.PlainCmdRunner{}.Run(cmd)
}

func (r ReleaseCmdRunner) RunWithCancel(cmd *goexec.Cmd, cancelCh chan struct{}) error {
	if r.fullOutput {
		cmd.Stdout = io.MultiWriter(r.log, cmd.Stdout)
		cmd.Stderr = io.MultiWriter(r.log, cmd.Stderr)
	}

	fmt.Fprintf(r.log, "==> Executing %s %v\n", cmd.Path, cmd.Args)
	defer fmt.Fprintf(r.log, "==> Finished executing %s\n\n", cmd.Path)

	return exec.PlainCmdRunner{}.RunWithCancel(cmd, cancelCh)
}
