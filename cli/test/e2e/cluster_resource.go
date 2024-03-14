// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type ClusterResource struct{}

func RemoveClusterResource(t *testing.T, kind, name, ns string, kubectl Kubectl) {
	_, err := kubectl.RunWithOpts([]string{"delete", kind, name, "-n", ns}, RunOpts{AllowError: true, NoNamespace: true})
	if err != nil {
		require.Contains(t, err.Error(), "Error from server (NotFound)", "Failed to delete resource %s/%s: %s")
	}
}
