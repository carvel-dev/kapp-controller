// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

// SandboxCmdRunner implements CmdRunner interface by running command
// in a new process/mount namespace.
type SandboxCmdRunner struct {
	local exec.CmdRunner
	opts  SandboxCmdRunnerOpts
}

// SandboxCmdRunnerOpts specifies sandboxing options.
type SandboxCmdRunnerOpts struct {
	// RequiresPosix indicates which commands should run in an environment
	// with an extended file system (e.g. include /lib, /lib64, /usr/*, ...).
	RequiresPosix map[string]bool
	// RequiresNetwork indicates which commands need network access.
	RequiresNetwork map[string]bool
}

var _ exec.CmdRunner = SandboxCmdRunner{}

// NewSandboxCmdRunner returns a new SandboxCmdRunner.
func NewSandboxCmdRunner(local exec.CmdRunner, opts SandboxCmdRunnerOpts) SandboxCmdRunner {
	return SandboxCmdRunner{local, opts}
}
