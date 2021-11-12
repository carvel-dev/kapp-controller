// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

func TestPackageInstalledDelete(t *testing.T) {
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
		kapp.Run([]string{"delete", "-a", appName})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("Adding and installing test package", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(yaml), AllowError: true,
		})
		require.NoError(t, err)

		_, err = kappCtrl.RunWithOpts([]string{"package", "installed", "create", "--package-install", pkgiName, "--package-name", packageMetadataName, "--version", packageVersion}, RunOpts{})
		require.NoError(t, err)
	})

	logger.Section("Delete installed package and created resources", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "installed", "delete", "--package-install", pkgiName}, RunOpts{})
		require.NoError(t, err)
	})

	logger.Section("Check if package is uninstalled", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "list", "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("Check for Service Account", func() {
		svcAccountName := fmt.Sprintf("%s-kapp-test-sa", pkgiName)
		_, err := kubectl.RunWithOpts([]string{"get", "sa", svcAccountName}, RunOpts{
			AllowError: true,
		})
		require.Contains(t, err.Error(), "not found")
	})

	logger.Section("Check for Cluster Role", func() {
		clusterRoleName := fmt.Sprintf("%s-kapp-test-cluster-role", pkgiName)
		_, err := kubectl.RunWithOpts([]string{"get", "clusterroles", clusterRoleName}, RunOpts{
			AllowError: true,
		})
		require.Contains(t, err.Error(), "not found")
	})

	logger.Section("Check for Cluster Role Binding", func() {
		clusterRoleBindingName := fmt.Sprintf("%s-kapp-test-cluster-rolebinding", pkgiName)
		_, err := kubectl.RunWithOpts([]string{"get", "clusterrolebindings", clusterRoleBindingName}, RunOpts{
			AllowError: true,
		})
		require.Contains(t, err.Error(), "not found")
	})

}
