// Copyright 2026 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"testing"
	"time"

	"carvel.dev/kapp-controller/test/e2e"
	"github.com/stretchr/testify/require"
)

func TestK8APIEndpoints(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kubectl := e2e.Kubectl{T: t, Namespace: env.Namespace, L: logger}

	logger.Section("aggregated APIService is Available", func() {
		require.Eventually(t, func() bool {
			out, err := kubectl.RunWithOpts([]string{
				"get", "apiservice", "v1alpha1.data.packaging.carvel.dev",
				"-o", "jsonpath={.status.conditions[?(@.type=='Available')].status}",
			}, e2e.RunOpts{NoNamespace: true, AllowError: true})
			return err == nil && out == "True"
		}, 60*time.Second, 2*time.Second, "APIService v1alpha1.data.packaging.carvel.dev never became Available")
	})

	logger.Section("data.packaging.carvel.dev/v1alpha1 discovery", func() {
		out, _ := kubectl.RunWithOpts([]string{"get", "--raw", "/apis/data.packaging.carvel.dev/v1alpha1"},
			e2e.RunOpts{NoNamespace: true})
		require.Contains(t, out, `"name":"packages"`)
		require.Contains(t, out, `"name":"packagemetadatas"`)
	})

	logger.Section("packaging.carvel.dev/v1alpha1 CRD discovery", func() {
		out, _ := kubectl.RunWithOpts([]string{"api-resources", "--api-group=packaging.carvel.dev", "-o", "name"},
			e2e.RunOpts{NoNamespace: true})
		require.Contains(t, out, "packageinstalls.packaging.carvel.dev")
		require.Contains(t, out, "packagerepositories.packaging.carvel.dev")
	})

	logger.Section("kappctrl.k14s.io/v1alpha1 CRD discovery", func() {
		out, _ := kubectl.RunWithOpts([]string{"api-resources", "--api-group=kappctrl.k14s.io", "-o", "name"},
			e2e.RunOpts{NoNamespace: true})
		require.Contains(t, out, "apps.kappctrl.k14s.io")
	})
}
