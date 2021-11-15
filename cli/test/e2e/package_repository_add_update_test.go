// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageRepositoryAddOrUpdate(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kapp{t, env.Namespace, env.KappCtrlBinaryPath, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	pkgrName := "test-package-repository"
	// TODO: Push repository todocker hub using k8slt
	pkgrURL := `index.docker.io/prewar/kc-cli-e2e-test-repo:latest`
	pkgrURLUpdated := `index.docker.io/k8slt/kc-e2e-test-repo:latest`

	kind := "PackageRepository"

	cleanUp := func() {
		RemoveClusterResource(t, kind, pkgrName, env.Namespace, kubectl)
	}

	cleanUp()
	defer cleanUp()

	logger.Section("package repository add", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "repository", "add", "-r", pkgrName, "--url", pkgrURL}, RunOpts{})
		require.NoError(t, err)

		kubectl.Run([]string{"get", kind, pkgrName})
		kubectl.Run([]string{"get", "pkgm/pkg.cli.test.carvel.dev"})
		kubectl.Run([]string{"get", "pkg/pkg.cli.test.carvel.dev.1.0.0"})
	})

	logger.Section("package repository update", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "repository", "update", "-r", pkgrName, "--url", pkgrURLUpdated}, RunOpts{})
		require.NoError(t, err)

		kubectl.Run([]string{"get", kind, pkgrName})
		kubectl.Run([]string{"get", "pkgm/pkg.test.carvel.dev"})
		kubectl.Run([]string{"get", "pkg/pkg.test.carvel.dev.1.0.0"})
		kubectl.Run([]string{"get", "pkg/pkg.test.carvel.dev.2.0.0"})
	})
}
