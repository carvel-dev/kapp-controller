// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"fmt"
	"io"
	goexec "os/exec"
	"path/filepath"
	"strings"

	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

type ReleaseCmdRunner struct {
	log                 io.Writer
	fullOutput          bool
	tempImgLockFilepath string
	ui                  cmdcore.AuthoringUI
}

var _ exec.CmdRunner = &ReleaseCmdRunner{}

func NewReleaseCmdRunner(log io.Writer, fullOutput bool, tempImgLockFilepath string, ui cmdcore.AuthoringUI) *ReleaseCmdRunner {
	return &ReleaseCmdRunner{log: log, fullOutput: fullOutput, tempImgLockFilepath: tempImgLockFilepath, ui: ui}
}

func (r ReleaseCmdRunner) Run(cmd *goexec.Cmd) error {
	if filepath.Base(cmd.Path) == "vendir" {
		r.ui.PrintInformationalText("kbld builds images when necessary and ensures that all image references are resolved to an immutable reference\n")
		r.ui.PrintHeaderText("Building images and resolving references")
	}

	if filepath.Base(cmd.Path) == "kapp" {
		return nil
	}
	if filepath.Base(cmd.Path) == "kbld" {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--imgpkg-lock-output=%s", r.tempImgLockFilepath))
	}
	if r.fullOutput {
		cmd.Stdout = io.MultiWriter(r.log, cmd.Stdout)
		cmd.Stderr = io.MultiWriter(r.log, cmd.Stderr)
	}

	if filepath.Base(cmd.Path) == "ytt" || filepath.Base(cmd.Path) == "kbld" {
		r.ui.PrintCmdExecutionOutput(fmt.Sprintf("$ %s", strings.Join(cmd.Args, " ")))
	}

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
