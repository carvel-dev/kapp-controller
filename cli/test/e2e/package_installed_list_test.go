// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

func TestPackageInstalledList(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	repoYml := `---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: e2e-repo.test.carvel.dev
spec:
  fetch:
    imgpkgBundle:
      image: index.docker.io/k8slt/kc-e2e-test-repo@sha256:ddd93b67b97c1460580ca1afd04326d16900dc716c4357cade85b83deab76f1c`

	pkgrName := "e2e-repo.test.carvel.dev"
	pkgName := "pkg.test.carvel.dev"
	pkgVersion := "1.0.0"
	pkgrKind := "PackageRepository"
	pkgiName := "testpkgi"

	cleanUp := func() {
		_, _ = kapp.RunWithOpts([]string{"package", "installed", "delete", "--name", pkgiName}, RunOpts{})
		RemoveClusterResource(t, pkgrKind, pkgrName, kapp.namespace, kubectl)
	}

	cleanUp()
	defer cleanUp()

	logger.Section("package installed list with no package present", func() {
		out, err := kapp.RunWithOpts([]string{"package", "available", "list", "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)

	})

	NewClusterResourceFromYaml(t, kubectl, pkgrKind, pkgrName, repoYml)
	_, err := kapp.RunWithOpts([]string{"package", "installed", "create", "--name", pkgiName, "--package-name", pkgName, "--version", pkgVersion}, RunOpts{})
	require.NoError(t, err)

	logger.Section("package installed list with one package installed", func() {
		out, err := kapp.RunWithOpts([]string{"package", "installed", "list", "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{
			{
				"name":            "testpkgi",
				"package_name":    "pkg.test.carvel.dev",
				"package_version": "1.0.0",
				"status":          "Reconcile succeeded",
			},
		}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})
}
