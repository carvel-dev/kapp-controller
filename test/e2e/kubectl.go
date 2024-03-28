// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

type Kubectl struct {
	T         *testing.T
	Namespace string
	L         Logger
}

func (k Kubectl) Run(args []string) string {
	out, _ := k.RunWithOpts(args, RunOpts{})
	return out
}

func (k Kubectl) RunWithOpts(args []string, opts RunOpts) (string, error) {
	if !opts.NoNamespace {
		args = append(args, []string{"-n", k.Namespace}...)
	}
	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	k.L.Debugf("Running '%s'...\n", k.cmdDesc(args))

	var stderr bytes.Buffer
	var stdout bytes.Buffer

	cmd := exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stderr = &stderr

	if opts.StdoutWriter != nil {
		cmd.Stdout = opts.StdoutWriter
	} else {
		cmd.Stdout = &stdout
	}

	cmd.Stdin = opts.StdinReader

	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("Execution error: stderr: '%s' error: '%s'", stderr.String(), err)

		if !opts.AllowError {
			k.T.Fatalf("Failed to successfully execute '%s': %v", k.cmdDesc(args), err)
		}
	}

	return stdout.String(), err
}

func (k Kubectl) cmdDesc(args []string) string {
	return fmt.Sprintf("kubectl %s", strings.Join(args, " "))
}
