// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	hasShowManagedFieldsFlag       bool
	determineShowManagedFieldsFlag sync.Once
)

type ClusterResource struct{}

func NewMissingClusterResource(t *testing.T, kind, name, ns string, kubectl Kubectl) {
	_, err := kubectl.RunWithOpts([]string{"get", kind, name, "-n", ns, "-o", "yaml"}, RunOpts{AllowError: true})
	if err == nil || !strings.Contains(err.Error(), "Error from server (NotFound)") {
		t.Fatalf("Expected resource to not exist")
	}
}

func NewClusterResourceFromYaml(t *testing.T, kubectl Kubectl, kind, name, yaml string) {
	_, err := kubectl.RunWithOpts([]string{"apply", "-f", "-"}, RunOpts{StdinReader: strings.NewReader(yaml)})
	require.NoError(t, err)

	start := time.Now()
	for {
		out, err := kubectl.RunWithOpts([]string{"get", kind, name}, RunOpts{})
		require.NoError(t, err)

		if strings.Contains(out, "Reconcile succeeded") {
			return
		}
		require.True(t, time.Now().Before(start.Add(25*time.Second)))
		time.Sleep(5 * time.Second)
	}
}

func RemoveClusterResource(t *testing.T, kind, name, ns string, kubectl Kubectl) {
	_, err := kubectl.RunWithOpts([]string{"delete", kind, name, "-n", ns}, RunOpts{AllowError: true, NoNamespace: true})
	if err != nil {
		require.Contains(t, err.Error(), "Error from server (NotFound)", "Failed to delete resource %s/%s: %s")
	}
}
