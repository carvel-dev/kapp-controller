// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"testing"
)

type Kapp struct {
	T         *testing.T
	Namespace string
	L         Logger
}

type RunOpts struct {
	NoNamespace  bool
	IntoNs       bool
	AllowError   bool
	StderrWriter io.Writer
	StdoutWriter io.Writer
	StdinReader  io.Reader
	Ctx          context.Context
	Redact       bool
	Interactive  bool
	OnErrKubectl []string
}

func (k Kapp) Run(args []string) string {
	out, _ := k.RunWithOpts(args, RunOpts{})
	return out
}

func (k Kapp) RunWithOpts(args []string, opts RunOpts) (string, error) {
	if !opts.NoNamespace {
		args = append(args, []string{"-n", k.Namespace}...)
	}
	if opts.IntoNs {
		args = append(args, []string{"--into-ns", k.Namespace}...)
	}
	if !opts.Interactive {
		args = append(args, "--yes")
	}
	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.TODO()
	}
	if args[0] == "deploy" {
		args = append(args, []string{"--wait-timeout", "3m"}...)
	}

	k.L.Debugf("Running '%s'...\n", k.cmdDesc(args, opts))

	cmdName := "kapp"
	cmd := exec.CommandContext(ctx, cmdName, args...)
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

	err := cmd.Run()
	stdoutStr := stdout.String()

	if err != nil {
		err = fmt.Errorf("Execution error: stdout: '%s' stderr: '%s' error: '%s'", stdoutStr, stderr.String(), err)

		if len(opts.OnErrKubectl) > 0 {
			kubectl := Kubectl{k.T, k.Namespace, k.L}
			debugOut, err := kubectl.RunWithOpts(opts.OnErrKubectl, RunOpts{AllowError: true})
			if err != nil {
				k.L.Debugf("OnErrKubectl error: %s\n", err)
			} else {
				k.L.Debugf("OnErrKubectl output:\n%s\n", debugOut)
			}
		}

		if !opts.AllowError {
			k.T.Fatalf("Failed to successfully execute '%s': %v", k.cmdDesc(args, opts), err)
		}
	}

	return stdoutStr, err
}

func (k Kapp) cmdDesc(args []string, opts RunOpts) string {
	prefix := "kapp"
	if opts.Redact {
		return prefix + " -redacted-"
	}
	return fmt.Sprintf("%s %s", prefix, strings.Join(args, " "))
}
