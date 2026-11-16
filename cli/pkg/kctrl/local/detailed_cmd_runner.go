// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package local

import (
	"fmt"
	"io"
	"os"
	goexec "os/exec"

	"carvel.dev/kapp-controller/pkg/exec"
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
	// Preserve child output for its consumer; never tee it into diagnostics.
	if r.fullOutput {
		_, _ = fmt.Fprintln(r.log, "Child output omitted from debug diagnostics")
	}

	// Adding os environment keys to cmd environment
	cmd.Env = append(os.Environ(), cmd.Env...)

	_, _ = fmt.Fprintf(r.log, "==> Executing %s\n", cmd.Path)
	defer fmt.Fprintf(r.log, "==> Finished executing %s\n\n", cmd.Path)

	return exec.PlainCmdRunner{}.Run(cmd)
}

func (r DetailedCmdRunner) RunWithCancel(cmd *goexec.Cmd, cancelCh chan struct{}) error {
	// Preserve child output for its consumer; never tee it into diagnostics.
	if r.fullOutput {
		_, _ = fmt.Fprintln(r.log, "Child output omitted from debug diagnostics")
	}

	_, _ = fmt.Fprintf(r.log, "==> Executing %s\n", cmd.Path)
	defer fmt.Fprintf(r.log, "==> Finished executing %s\n\n", cmd.Path)

	return exec.PlainCmdRunner{}.RunWithCancel(cmd, cancelCh)
}
