// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"os"
	"strings"
	"testing"
)

type Env struct {
	Namespace          string
	KappBinaryPath     string
	KappCtrlBinaryPath string
}

func BuildEnv(t *testing.T) Env {
	kappPath := os.Getenv("KAPP_BINARY_PATH")
	if kappPath == "" {
		kappPath = "kapp"
	}

	kappCtrlPath := os.Getenv("KAPPCTRL_BINARY_PATH")
	if kappCtrlPath == "" {
		kappCtrlPath = "kapp"
	}

	env := Env{
		Namespace:          os.Getenv("KAPP_E2E_NAMESPACE"),
		KappBinaryPath:     kappPath,
		KappCtrlBinaryPath: kappCtrlPath,
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
