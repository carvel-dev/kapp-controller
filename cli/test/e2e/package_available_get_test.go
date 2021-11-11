// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

func TestPackageAvailableGet(t *testing.T) {
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
	pkgrKind := "PackageRepository"

	cleanUp := func() {
		RemoveClusterResource(t, pkgrKind, pkgrName, kapp.namespace, kubectl)
	}

	cleanUp()
	defer cleanUp()

	NewClusterResourceFromYaml(t, kubectl, pkgrKind, pkgrName, repoYml)

	logger.Section("package available get", func() {
		out, err := kapp.RunWithOpts([]string{"package", "available", "get", "-p", pkgName, "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{
			{
				"category":          "testing",
				"display_name":      "Carvel Test Package",
				"long_description":  "This is a test application which has been packaged using the Carvel tools and can be deployed with them. For more information you can visit https://carvel.dev/.",
				"maintainers":       "- name: Carvel Team",
				"name":              "pkg.test.carvel.dev",
				"package_provider":  "Carvel",
				"short_description": "Carvel package for testing installation",
				"support":           "Visit #carvel in slack: https://kubernetes.slack.com/archives/CH8KCCKA5.",
			},
		}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})
}
