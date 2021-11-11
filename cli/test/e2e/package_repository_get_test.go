// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

func TestPackageRepositoryGet(t *testing.T) {
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

	NewClusterResourceFromYaml(t, kubectl, pkgrKind, pkgrName, repoYml)

	logger.Section("package repository get", func() {
		out, err := kapp.RunWithOpts([]string{"package", "repository", "get", "-r", pkgrName, "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{
			{
				"name":       "e2e-repo.test.carvel.dev",
				"reason":     "",
				"repository": "index.docker.io/k8slt/kc-e2e-test-repo",
				"status":     "Reconcile succeeded",
				"tag":        "sha256:ddd93b67b97c1460580ca1afd04326d16900dc716c4357cade85b83deab76f1c",
				"version":    "<replaced>",
			},
		}

		require.Exactly(t, expectedOutputRows, replaceVersion(output.Tables[0].Rows))
	})
}

func replaceVersion(result []map[string]string) []map[string]string {
	for i, row := range result {
		if len(row["version"]) > 0 {
			row["version"] = "<replaced>"
		}
		result[i] = row
	}
	return result
}
