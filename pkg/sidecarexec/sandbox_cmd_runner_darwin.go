// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	goexec "os/exec"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

func (SandboxCmdRunner) Run(cmd *goexec.Cmd, opts exec.RunOpts) error {
	panic("Not implemented on this OS")
}

func (SandboxCmdRunner) RunWithCancel(cmd *goexec.Cmd, cancelCh chan struct{}, opts exec.RunOpts) error {
	panic("Not implemented on this OS")
}

type SandboxWrap struct{}

func (SandboxWrap) ExecuteCmd(rawArgs []string) error {
	panic("Not implemented on this OS")
}
