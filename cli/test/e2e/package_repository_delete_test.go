// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageRepositoryDelete(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kapp{t, env.Namespace, env.KappCtrlBinaryPath, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	appName := "test-package-repository-delete"

	repoYml := `---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: e2e-repo.test.carvel.dev
spec:
  fetch:
    imgpkgBundle:
      image: index.docker.io/k8slt/kc-e2e-test-repo@sha256:ddd93b67b97c1460580ca1afd04326d16900dc716c4357cade85b83deab76f1c`

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName})
	}

	cleanUp()
	defer cleanUp()

	pkgrName := "e2e-repo.test.carvel.dev"

	logger.Section("deletion of nonexistent repository", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "repository", "delete", "-r", pkgrName}, RunOpts{
			AllowError: true,
		})
		require.Error(t, err)
	})

	logger.Section("deletion of an existing repository", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(repoYml),
		})

		kubectl.Run([]string{"get", "packagerepository", pkgrName})

		kappCtrl.RunWithOpts([]string{"package", "repository", "delete", "-r", pkgrName}, RunOpts{})

		_, err := kubectl.RunWithOpts([]string{"get", "packagerepository", pkgrName}, RunOpts{AllowError: true})
		require.Contains(t, err.Error(), "not found")
	})
}
