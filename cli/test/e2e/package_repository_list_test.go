// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

func TestPackageRepositoryList(t *testing.T) {
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
	pkgrKind := "PackageRepository"

	cleanUp := func() {
		RemoveClusterResource(t, pkgrKind, pkgrName, kapp.namespace, kubectl)
	}

	cleanUp()
	defer cleanUp()

	logger.Section("package repository list with no package present", func() {
		out, err := kapp.RunWithOpts([]string{"package", "repository", "list", "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)

	})

	NewClusterResourceFromYaml(t, kubectl, pkgrKind, pkgrName, repoYml)

	logger.Section("package repository list with one package installed", func() {
		out, err := kapp.RunWithOpts([]string{"package", "installed", "list", "--json", "-A"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		//TODO: Update once we handle URLs better
		expectedOutputRows := []map[string]string{
			{
				"name":            "testpkgi",
				"namespace":       "default",
				"package_name":    "pkg.test.carvel.dev",
				"package_version": "1.0.0",
				"status":          "Reconcile succeeded",
			},
		}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})
}
