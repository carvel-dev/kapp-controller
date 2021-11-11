// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageInstalledCreate(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kapp{t, env.Namespace, env.KappCtrlBinaryPath, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	appName := "test-package-name"
	pkgiName := "testpkgi"
	packageMetadataName := "test-pkg.carvel.dev"

	packageMetadata := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: %s
spec:
  displayName: "Carvel Test Package"
  shortDescription: "Carvel package for testing installation"`, packageMetadataName)

	packageName := "test-pkg.carvel.dev.1.0.0"
	packageVersion := "1.0.0"

	packageCR := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: %s
spec:
  refName: test-pkg.carvel.dev
  version: %s
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: k8slt/kctrl-example-pkg:v1.0.0
      template:
      - ytt:
          paths:
          - config/
      - kbld:
          paths:
          - "-"
          - ".imgpkg/images.yml"
      deploy:
      - kapp: {}`, packageName, packageVersion)

	yaml := packageMetadata + "\n" + packageCR

	//TODO: Add test with data values. Let --values-file

	cleanUp := func() {
		// TODO: Check for error while uninstalling in cleanup?
		kappCtrl.Run([]string{"package", "installed", "delete", "--name", pkgiName})
		kapp.Run([]string{"delete", "-a", appName})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("Adding test package", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(yaml), AllowError: true,
		})
		require.NoError(t, err)
	})

	logger.Section("Installing test package", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "installed", "create", "--name", pkgiName, "--package-name", packageMetadataName, "--version", packageVersion}, RunOpts{})
		require.NoError(t, err)
	})

	svcAccountName := fmt.Sprintf("%s-kapp-test-sa", pkgiName)
	logger.Section("Check for Service Account", func() {
		_, err := kubectl.RunWithOpts([]string{"get", "sa", svcAccountName}, RunOpts{})
		require.NoError(t, err)
	})

	clusterRoleName := fmt.Sprintf("%s-kapp-test-cluster-role", pkgiName)
	logger.Section("Check for Cluster Role", func() {
		_, err := kubectl.RunWithOpts([]string{"get", "clusterroles", clusterRoleName}, RunOpts{})
		require.NoError(t, err)
	})

	clusterRoleBindingName := fmt.Sprintf("%s-kapp-test-cluster-rolebinding", pkgiName)
	logger.Section("Check for Cluster Role Binding", func() {
		_, err := kubectl.RunWithOpts([]string{"get", "clusterrolebindings", clusterRoleBindingName}, RunOpts{})
		require.NoError(t, err)
	})

}
