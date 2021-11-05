// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageAvailableList(t *testing.T) {
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

	name := "e2e-repo.test.carvel.dev"
	kind := "PackageRepository"
	cleanUp := func() {
		RemoveClusterResource(t, kind, name, kapp.namespace, kubectl)
	}

	cleanUp()
	defer cleanUp()

	logger.Section("package available list with no package present", func() {
		out, err := kapp.RunWithOpts([]string{"package", "available", "list"}, RunOpts{})
		require.NoError(t, err)

		out = strings.TrimSpace(replaceSpaces(replaceTarget(out)))

		expectedOutput := strings.TrimSpace(replaceSpaces(`

Available packages in namespace 'kapp-test'

Name  Display-Name  Short-Description  

0 Packages Available

Succeeded`))

		require.Equal(t, expectedOutput, out)

	})

	NewClusterResourceFromYaml(t, kubectl, kind, name, repoYml)

	logger.Section("package available list with one package available", func() {
		out, err := kapp.RunWithOpts([]string{"package", "available", "list"}, RunOpts{})
		require.NoError(t, err)

		out = strings.TrimSpace(replaceSpaces(replaceTarget(out)))

		expectedOutput := strings.TrimSpace(replaceSpaces(`

Available packages in namespace 'kapp-test'

Name                 Display-Name         Short-Description  
pkg.test.carvel.dev  Carvel Test Package  Carvel package for testing installation  

1 Packages Available

Succeeded`))

		require.Equal(t, expectedOutput, out)
	})

	logger.Section("package available list versions of a package", func() {
		out, err := kapp.RunWithOpts([]string{"package", "available", "list", "-P", "pkg.test.carvel.dev"}, RunOpts{})
		require.NoError(t, err)

		out = strings.TrimSpace(replaceSpaces(replaceTarget(out)))

		expectedOutput := strings.TrimSpace(replaceSpaces(`

Available package versions for 'pkg.test.carvel.dev' in namespace 'kapp-test'

Name                 Version     Released-At  
pkg.test.carvel.dev  1.0.0       0001-01-01 00:00:00 +0000 UTC  
pkg.test.carvel.dev  2.0.0       0001-01-01 00:00:00 +0000 UTC  
pkg.test.carvel.dev  3.0.0-rc.1  0001-01-01 00:00:00 +0000 UTC  

3 Package Versions Available

Succeeded`))

		require.Equal(t, expectedOutput, out)
	})
}

func replaceSpaces(result string) string {
	// result = strings.Replace(result, " ", "_", -1) // useful for debugging
	result = strings.Replace(result, " \n", " $\n", -1) // explicit endline
	return result
}

func replaceTarget(result string) string {
	return regexp.MustCompile("Target cluster .+\n").ReplaceAllString(result, "")
}
