// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
)

// TestSAandClusterCELValidation tests packageInstall and App CR
// to ensure that either SA or cluster is validated at admission.
func TestSAandClusterCELValidation(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

  // Check the Kubernetes version and skip the test if it's less than 1.26
	kubeVersion, err := env.GetServerVersion()
	if err != nil || kubeVersion == nil {
		t.Fatalf("Error getting Kubernetes version: %v", err)
	}

	if kubeVersion.Major < "1" || (kubeVersion.Major == "1" && kubeVersion.Minor < "26") {
		t.Skip("Skipping test for Kubernetes versions less than 1.26")
	}

	name := "incorrect-spec-without-sa-cluster"

	appYAML := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  fetch:
    - inline:
        paths:
          file.yml: |
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
  template:
  - ytt: {}
  deploy:
    - kapp: {}
`, name)

	pkginstallYAML := fmt.Sprintf(`
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %[2]s
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
spec:
  packageRef:
    refName: pkg.incorrect.carvel.dev
    versionSelection:
      constraints: 1.0.0
`, env.Namespace, name)

	logger.Section("Create App CR with kubectl", func() {
		_, err := kubectl.RunWithOpts([]string{"apply", "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(appYAML), AllowError: true})
		require.Error(t, err)
		require.ErrorContains(t, err, "Expected service account or cluster.")
	})

	logger.Section("Create PackageInstall with kapp", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(pkginstallYAML), AllowError: true})
		require.Error(t, err)
		require.ErrorContains(t, err, "Expected service account or cluster.")
	})
}