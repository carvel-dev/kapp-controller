// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"strings"
	"testing"
)

func TestHelpCommandGroup(t *testing.T) {
	env := BuildEnv(t)
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, Logger{}}

	_, err := kappCtrl.RunWithOpts([]string{"package"}, RunOpts{NoNamespace: true, AllowError: true})
	if err == nil {
		t.Fatalf("Expected error")
	}
	if !strings.Contains(err.Error(), "Error: Use one of available subcommands: available, init, install, installed, release, repository") {
		t.Fatalf("Expected helpful error message, but was '%s'", err.Error())
	}
}
