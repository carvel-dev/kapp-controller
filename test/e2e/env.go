// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"os"
	"strings"
	"testing"
)

type Env struct {
	Namespace         string
	PackagingGlobalNS string
}

func BuildEnv(t *testing.T) Env {
	env := Env{
		Namespace: os.Getenv("KAPPCTRL_E2E_NAMESPACE"),
	}

	if pkgNS := os.Getenv("KAPPCTRL_E2E_PACKAGING_NAMESPACE"); pkgNS != "" {
		env.PackagingGlobalNS = pkgNS
	} else {
		env.PackagingGlobalNS = "kapp-controller-packaging-global"
	}

	env.Validate(t)
	return env
}

func (e Env) Validate(t *testing.T) {
	errStrs := []string{}

	if len(e.Namespace) == 0 {
		errStrs = append(errStrs, "Expected Namespace to be non-empty")
	}

	if len(errStrs) > 0 {
		t.Fatalf("%s", strings.Join(errStrs, "\n"))
	}
}
