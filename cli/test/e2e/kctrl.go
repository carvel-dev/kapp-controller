// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type Kctrl struct {
	t         *testing.T
	namespace string
	kctrlPath string
	l         Logger
}

func (k Kctrl) Run(args []string) string {
	out, _ := k.RunWithOpts(args, RunOpts{})
	return out
}

func (k Kctrl) RunWithOpts(args []string, opts RunOpts) (string, error) {
	if !opts.NoNamespace {
		args = append(args, []string{"-n", k.namespace}...)
	}
	if opts.IntoNs {
		args = append(args, []string{"--into-ns", k.namespace}...)
	}
	if !opts.Interactive {
		args = append(args, "--yes")
	}

	k.l.Debugf("Running '%s'...\n", k.cmdDesc(args, opts))

	cmd := exec.Command(k.kctrlPath, args...)
	cmd.Stdin = opts.StdinReader

	var stderr, stdout bytes.Buffer

	if opts.StderrWriter != nil {
		cmd.Stderr = opts.StderrWriter
	} else {
		cmd.Stderr = &stderr
	}

	if opts.StdoutWriter != nil {
		cmd.Stdout = opts.StdoutWriter
	} else {
		cmd.Stdout = &stdout
	}

	if opts.CancelCh != nil {
		go func() {
			select {
			case <-opts.CancelCh:
				cmd.Process.Signal(os.Interrupt)
			}
		}()
	}

	err := cmd.Run()
	stdoutStr := stdout.String()

	if err != nil {
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}

		err = fmt.Errorf("Execution error: stdout: '%s' stderr: '%s' error: '%s' exit code: '%d'",
			stdoutStr, stderr.String(), err, exitCode)

		if !opts.AllowError {
			k.t.Fatalf("Failed to successfully execute '%s': %v", k.cmdDesc(args, opts), err)
		}
	}

	return stdoutStr, err
}

func (k Kctrl) cmdDesc(args []string, opts RunOpts) string {
	prefix := "kctrl"
	if opts.Redact {
		return prefix + " -redacted-"
	}
	return fmt.Sprintf("%s %s", prefix, strings.Join(args, " "))
}
