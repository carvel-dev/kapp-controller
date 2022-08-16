// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"os"
	"strings"
	"testing"
)

type Env struct {
	Namespace       string
	Image           string
	KappBinaryPath  string
	KctrlBinaryPath string
}

func BuildEnv(t *testing.T) Env {
	kappPath := os.Getenv("KAPP_BINARY_PATH")
	if kappPath == "" {
		kappPath = "kapp"
	}

	kctrlPath := os.Getenv("KCTRL_BINARY_PATH")
	if kctrlPath == "" {
		kctrlPath = "kctrl"
	}

	env := Env{
		Namespace:       os.Getenv("KCTRL_E2E_NAMESPACE"),
		Image:           os.Getenv("KCTRL_E2E_IMAGE"),
		KappBinaryPath:  kappPath,
		KctrlBinaryPath: kctrlPath,
	}
	env.Validate(t)
	return env
}

func (e Env) Validate(t *testing.T) {
	errStrs := []string{}

	if len(e.Namespace) == 0 {
		errStrs = append(errStrs, "Expected Namespace to be non-empty (hint: kubectl create namespace kctrl-test; export KCTRL_E2E_NAMESPACE=kctrl-test)")
	}

	if len(e.Image) == 0 {
		errStrs = append(errStrs, "Test image to be non-empty (hint: export KCTRL_E2E_IMAGE=repo/kctrl-test)")
	}

	if len(errStrs) > 0 {
		t.Fatalf("%s", strings.Join(errStrs, "\n"))
	}
}
