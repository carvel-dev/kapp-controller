// Copyright 2026 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"errors"
	"strings"
	"testing"
)

func TestResultRedactsAllDiagnosticFields(t *testing.T) {
	raw := "https://user:password-canary@example.com/path?access_token=query-canary#fragment-canary"
	original := CmdRunResult{Stdout: raw, Stderr: raw, Error: errors.New(raw), ExitCode: 1, Finished: true}
	safe := original.WithRedactedSecrets()
	if strings.Contains(safe.Stdout+safe.Stderr+safe.ErrorStr(), "canary") {
		t.Fatal("credential in diagnostics")
	}
	if original.Stdout != raw || original.ErrorStr() != raw || safe.ExitCode != 1 || !safe.Finished {
		t.Fatal("execution result changed")
	}
}
