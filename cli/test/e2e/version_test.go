// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	env := BuildEnv(t)
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, Logger{}}

	out, _ := kappCtrl.RunWithOpts([]string{"version"}, RunOpts{NoNamespace: true})

	if !strings.Contains(out, "kctrl version") {
		t.Fatalf("Expected to find client version")
	}
}
