// Copyright 2026 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package local

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestDetailedRunnerDoesNotTeeChildOutput(t *testing.T) {
	for _, debug := range []bool{false, true} {
		for _, cancelable := range []bool{false, true} {
			var log, stdout, stderr bytes.Buffer
			cmd := exec.Command("sh", "-c", `printf '%s' "$SENSITIVE_VALUE"; printf '%s' "$SENSITIVE_VALUE" >&2`)
			cmd.Env = append(cmd.Environ(), "SENSITIVE_VALUE=password-canary")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			runner := NewDetailedCmdRunner(&log, debug)
			var err error
			if cancelable {
				err = runner.RunWithCancel(cmd, make(chan struct{}))
			} else {
				err = runner.Run(cmd)
			}
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(log.String(), "canary") || stdout.String() != "password-canary" || stderr.String() != "password-canary" {
				t.Fatal("unsafe diagnostics or changed child output")
			}
		}
	}
}
